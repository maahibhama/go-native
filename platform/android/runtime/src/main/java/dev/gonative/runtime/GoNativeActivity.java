package dev.gonative.runtime;

import android.app.Activity;
import android.os.Bundle;
import android.os.Looper;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.view.WindowInsets;
import android.widget.HorizontalScrollView;
import android.widget.ScrollView;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicReference;

public class GoNativeActivity extends Activity implements EventDispatcher {
    private static final int CREATE = 1;
    private static final int DELETE = 2;
    private static final int UPDATE = 3;
    private static final int INSERT = 4;
    private static final int REMOVE = 5;
    private static final int MOVE = 6;

    static { System.loadLibrary("gonative"); }

    private final ViewRegistry registry = new ViewRegistry();
    private ControlFactory controlFactory;
    private NativeMeasurer measurer;
    private float lastViewportWidth, lastViewportHeight, lastViewportScale;
    private final View.OnLayoutChangeListener viewportListener = new View.OnLayoutChangeListener() {
        @Override public void onLayoutChange(View view, int left, int top, int right, int bottom,
                                             int oldLeft, int oldTop, int oldRight, int oldBottom) {
            reportViewport(view);
        }
    };

    private native void nativeStart();
    private native void nativeDispatchEvent(long handler);
    private native void nativeDispatchValueEvent(long handler, String value);
    private native void nativeDispatchBoolEvent(long handler, boolean value);
    private native void nativeDispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY);
    private native void nativeDispatchSelectionEvent(long handler, int start, int end);
    private native void nativeStop();
    private native void nativeSetLifecycle(int state);
    private native void nativeDispatchFocus(long nodeID, boolean focused);
    private native void nativeUpdateViewport(float width, float height, float scale);
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override public void dispatchEvent(long handler) { nativeDispatchEvent(handler); }
    @Override public void dispatchValueEvent(long handler, String value) { nativeDispatchValueEvent(handler, value); }
    @Override public void dispatchBoolEvent(long handler, boolean value) { nativeDispatchBoolEvent(handler, value); }
    @Override public void dispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY) {
        nativeDispatchGestureEvent(handler, translationX, translationY, velocityX, velocityY);
    }
    @Override public void dispatchSelectionEvent(long handler, int start, int end) { nativeDispatchSelectionEvent(handler, start, end); }
    @Override public void dispatchFocus(long nodeID, boolean focused) { nativeDispatchFocus(nodeID, focused); }

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        controlFactory = new ControlFactory(this, registry, this);
        measurer = new NativeMeasurer(this, controlFactory);
        getWindow().getDecorView().setBackgroundColor(android.graphics.Color.WHITE);
        nativeStart();
        nativeSetLifecycle(0);
        View content = findViewById(android.R.id.content);
        content.setOnApplyWindowInsetsListener(new View.OnApplyWindowInsetsListener() {
            @Override public WindowInsets onApplyWindowInsets(View view, WindowInsets insets) {
                int left = insets.getSystemWindowInsetLeft();
                int top = insets.getSystemWindowInsetTop();
                int right = insets.getSystemWindowInsetRight();
                int bottom = insets.getSystemWindowInsetBottom();
                view.setPadding(left, top, right, bottom);
                reportViewport(view);
                return insets;
            }
        });
        content.addOnLayoutChangeListener(viewportListener);
        content.requestApplyInsets();
        content.post(new Runnable() { @Override public void run() { reportViewport(findViewById(android.R.id.content)); } });
    }

    @Override protected void onStart() { super.onStart(); nativeSetLifecycle(1); }
    @Override protected void onResume() { super.onResume(); nativeSetLifecycle(2); }
    @Override protected void onPause() { nativeSetLifecycle(3); super.onPause(); }
    @Override protected void onStop() { nativeSetLifecycle(4); super.onStop(); }
    @Override public void onLowMemory() { nativeSetLifecycle(5); super.onLowMemory(); }

    @Override protected void onDestroy() {
        nativeSetLifecycle(6);
        View content = findViewById(android.R.id.content);
        if (content != null) content.removeOnLayoutChangeListener(viewportListener);
        nativeStop();
        registry.clear();
        super.onDestroy();
    }

    private void reportViewport(View content) {
        if (content == null || content.getWidth() <= 0 || content.getHeight() <= 0) return;
        float scale = getResources().getDisplayMetrics().density;
        float width = (content.getWidth() - content.getPaddingLeft() - content.getPaddingRight()) / scale;
        float height = (content.getHeight() - content.getPaddingTop() - content.getPaddingBottom()) / scale;
        if (width == lastViewportWidth && height == lastViewportHeight && scale == lastViewportScale) return;
        lastViewportWidth = width; lastViewportHeight = height; lastViewportScale = scale;
        nativeUpdateViewport(width, height, scale);
    }

    @SuppressWarnings("unused")
    public byte[] measureNativeBatch(final byte[] payload) {
        if (payload == null || payload.length == 0 || payload.length > 16777216) return null;
        if (measurer == null) {
            if (controlFactory == null) controlFactory = new ControlFactory(this, registry, this);
            measurer = new NativeMeasurer(this, controlFactory);
        }
        if (Looper.myLooper() == Looper.getMainLooper()) return measurer.measureNativeBatchOnUiThread(payload.clone());
        final AtomicReference<byte[]> result = new AtomicReference<>();
        final CountDownLatch ready = new CountDownLatch(1);
        runOnUiThread(new Runnable() {
            @Override public void run() {
                try { result.set(measurer.measureNativeBatchOnUiThread(payload.clone())); }
                finally { ready.countDown(); }
            }
        });
        try {
            return ready.await(10, TimeUnit.SECONDS) ? result.get() : null;
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
            return null;
        }
    }

    @SuppressWarnings("unused")
    public void applyMutationBatch(byte[] payload) {
        final byte[] owned = payload.clone();
        runOnUiThread(new Runnable() {
            @Override public void run() { applyOnUiThread(owned); }
        });
    }

    private void applyOnUiThread(byte[] payload) {
        if (payload == null) return;
        if (controlFactory == null) controlFactory = new ControlFactory(this, registry, this);
        try {
            long started = System.nanoTime();
            ProtocolReader in = new ProtocolReader(payload);
            if (in.remaining() < 14) return;
            if (in.uint16() != 10) return;
            int count = in.int32();
            if (count < 0 || count > 100000) return;
            long sequence = in.int64();
            for (int operation = 0; operation < count && in.hasRemaining(); operation++) {
                int mutation = in.uint8();
                int kind = in.uint8();
                long nodeID = in.int64();
                long parentID = in.int64();
                int index = in.int32();
                int fromIndex = in.int32();
                float width = in.float32();
                float height = in.float32();
                float padding = in.float32();
                float gap = in.float32();
                int alignment = in.uint8();
                boolean bold = in.uint8() != 0;
                float fontSize = in.float32();
                long handler = in.int64();
                long changeHandler = in.int64();
                long toggleHandler = in.int64();
                boolean checked = in.uint8() != 0;
                float progress = in.float32();
                String text = in.requiredString();
                String accessibility = in.requiredString();
                String hint = in.requiredString();
                int role = in.uint8();
                boolean focused = in.uint8() != 0;
                boolean scalesText = in.uint8() != 0;
                String imageSource = in.requiredString();
                int imageMode = in.uint8();
                boolean horizontal = in.uint8() != 0;
                int interactionLength = in.remaining() >= 4 ? in.int32() : 0;
                if (interactionLength < 0 || interactionLength > 1048576 || interactionLength > in.remaining()) return;
                byte[] interactions = in.bytes(interactionLength);
                if (in.remaining() < 9) return;
                int textWrap = in.uint8();
                int textOverflow = in.uint8();
                long maxLinesValue = in.uint32();
                if (maxLinesValue > Integer.MAX_VALUE) return;
                int maxLines = (int) maxLinesValue;
                boolean selectable = in.uint8() != 0;
                int richTextLength = in.int32();
                if (richTextLength < 0 || richTextLength > 1048576 || richTextLength > in.remaining()) return;
                byte[] richText = in.bytes(richTextLength);
                if (in.remaining() < 12) return;
                long linkHandler = in.int64();
                String placeholder = in.requiredString();
                if (in.remaining() < 36) return;
                int inputMode = in.uint8();
                int inputKind = in.uint8();
                int returnKey = in.uint8();
                int capitalization = in.uint8();
                int autoCorrect = in.uint8();
                boolean secure = in.uint8() != 0;
                boolean multiline = in.uint8() != 0;
                boolean readOnly = in.uint8() != 0;
                int validationState = in.uint8();
                String errorText = in.requiredString();
                if (in.remaining() < 32) return;
                int selectionStart = in.int32();
                int selectionEnd = in.int32();
                int maxLength = in.int32();
                long submitHandler = in.int64();
                long selectionHandler = in.int64();
                int styleLength = in.remaining() >= 4 ? in.int32() : -1;
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return;
                byte[] typedStyle = in.bytes(styleLength);
                if (in.remaining() < 17) return;
                boolean hasFrame = in.uint8() != 0;
                float frameX = in.float32(), frameY = in.float32(), frameWidth = in.float32(), frameHeight = in.float32();
                View view = registry.getView(nodeID);

                if (mutation == CREATE) {
                    view = controlFactory.makeView(kind, horizontal);
                    view.setTag(nodeID);
                    final long focusNodeID = nodeID;
                    view.setOnFocusChangeListener(new View.OnFocusChangeListener() {
                        @Override public void onFocusChange(View changed, boolean hasFocus) { nativeDispatchFocus(focusNodeID, hasFocus); }
                    });
                    registry.putView(nodeID, view);
                    if (registry.getRootNodeID() == 0) registry.setRootNodeID(nodeID);
                    controlFactory.style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode, textWrap, textOverflow, maxLines, selectable, richText, linkHandler, placeholder, inputMode, inputKind, returnKey, capitalization, autoCorrect, secure, multiline, readOnly, validationState, errorText, selectionStart, selectionEnd, maxLength, submitHandler, selectionHandler);
                    controlFactory.applyTypedStyle(view, typedStyle);
                    controlFactory.applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight, nodeID == registry.getRootNodeID());
                    controlFactory.applyInteractions(nodeID, view, interactions);
                    if (registry.viewCount() == 1) {
                        view.setLayoutParams(new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
                        view.setBackgroundColor(android.graphics.Color.WHITE);
                        setContentView(view);
                    }
                } else if (mutation == UPDATE) {
                    if (view != null) {
                        controlFactory.style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode, textWrap, textOverflow, maxLines, selectable, richText, linkHandler, placeholder, inputMode, inputKind, returnKey, capitalization, autoCorrect, secure, multiline, readOnly, validationState, errorText, selectionStart, selectionEnd, maxLength, submitHandler, selectionHandler);
                        controlFactory.applyTypedStyle(view, typedStyle);
                        controlFactory.applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight, nodeID == registry.getRootNodeID());
                        controlFactory.applyInteractions(nodeID, view, interactions);
                    }
                } else if (mutation == INSERT) {
                    View parentView = registry.getView(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        ViewRegistry.detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == REMOVE) {
                    ViewRegistry.detach(view);
                } else if (mutation == MOVE) {
                    View parentView = registry.getView(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        ViewRegistry.detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == DELETE) {
                    registry.deleteNode(nodeID);
                }
                if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
            }
            nativeReportBatchApplied(sequence, System.nanoTime() - started);
        } catch (Throwable t) {
            Log.e("GoNative", "Error applying mutation batch", t);
        }
    }
}
