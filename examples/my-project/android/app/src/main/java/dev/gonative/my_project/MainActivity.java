package dev.gonative.my_project;

import android.app.Activity;
import android.graphics.Typeface;
import android.os.Bundle;
import android.os.Looper;
import android.util.LongSparseArray;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.view.accessibility.AccessibilityEvent;
import android.view.accessibility.AccessibilityNodeInfo;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;
import android.widget.Switch;
import android.widget.ProgressBar;
import android.widget.CompoundButton;
import android.widget.ImageView;
import android.widget.ScrollView;
import android.widget.HorizontalScrollView;
import android.text.Editable;
import android.text.TextWatcher;
import android.animation.Animator;
import android.animation.AnimatorSet;
import android.animation.ObjectAnimator;
import android.animation.ValueAnimator;
import android.view.MotionEvent;
import android.view.VelocityTracker;
import android.view.animation.AccelerateDecelerateInterpolator;
import android.view.animation.AccelerateInterpolator;
import android.view.animation.DecelerateInterpolator;
import android.view.animation.LinearInterpolator;
import android.view.animation.OvershootInterpolator;

import java.util.ArrayList;
import java.util.List;
import java.util.WeakHashMap;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicReference;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;

public final class MainActivity extends Activity {
    private static final int CREATE = 1;
    private static final int DELETE = 2;
    private static final int UPDATE = 3;
    private static final int INSERT = 4;
    private static final int REMOVE = 5;
    private static final int MOVE = 6;
    private static final int TEXT = 2;
    private static final int BUTTON = 3;
    private static final int ROW = 4;
    private static final int COLUMN = 5;
    private static final int SAFE_AREA = 6;
    private static final int TEXT_INPUT = 7;
    private static final int SWITCH = 8;
    private static final int PROGRESS_INDICATOR = 9;
    private static final int IMAGE = 10;
    private static final int SCROLL_VIEW = 11;

    static { System.loadLibrary("gonative"); }

    private final LongSparseArray<View> views = new LongSparseArray<>();
    private final LongSparseArray<GestureBinding> gestureBindings = new LongSparseArray<>();
    private final WeakHashMap<EditText, TextWatcher> textWatchers = new WeakHashMap<>();
    private long rootNodeID;
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
    private native void nativeStop();
    private native void nativeSetLifecycle(int state);
    private native void nativeDispatchFocus(long nodeID, boolean focused);
    private native void nativeUpdateViewport(float width, float height, float scale);
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        getWindow().getDecorView().setBackgroundColor(android.graphics.Color.WHITE);
        nativeStart();
        nativeSetLifecycle(0);
        View content = findViewById(android.R.id.content);
        content.addOnLayoutChangeListener(viewportListener);
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
        for (int i = 0; i < gestureBindings.size(); i++) gestureBindings.valueAt(i).dispose();
        gestureBindings.clear();
        textWatchers.clear();
        for (int i = 0; i < views.size(); i++) views.valueAt(i).animate().cancel();
        views.clear();
        super.onDestroy();
    }

    private void reportViewport(View content) {
        if (content == null || content.getWidth() <= 0 || content.getHeight() <= 0) return;
        float scale = getResources().getDisplayMetrics().density;
        float width = content.getWidth() / scale, height = content.getHeight() / scale;
        if (width == lastViewportWidth && height == lastViewportHeight && scale == lastViewportScale) return;
        lastViewportWidth = width; lastViewportHeight = height; lastViewportScale = scale;
        nativeUpdateViewport(width, height, scale);
    }

    @SuppressWarnings("unused")
    public byte[] measureNativeBatch(final byte[] payload) {
        if (payload == null || payload.length == 0 || payload.length > 16777216) return null;
        if (Looper.myLooper() == Looper.getMainLooper()) return measureNativeBatchOnUiThread(payload.clone());
        final AtomicReference<byte[]> result = new AtomicReference<>();
        final CountDownLatch ready = new CountDownLatch(1);
        runOnUiThread(new Runnable() {
            @Override public void run() {
                try { result.set(measureNativeBatchOnUiThread(payload.clone())); }
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

    private byte[] measureNativeBatchOnUiThread(byte[] payload) {
        try {
            ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
            if (in.remaining() < 6 || Short.toUnsignedInt(in.getShort()) != 1) return null;
            int count = in.getInt();
            if (count < 0 || count > 100000) return null;
            ArrayList<NativeMeasurement> measured = new ArrayList<>(count);
            for (int i = 0; i < count; i++) {
                if (in.remaining() < 25) return null;
                long id = in.getLong();
                int kind = Byte.toUnsignedInt(in.get());
                float minWidth = in.getFloat(), maxWidth = in.getFloat(), minHeight = in.getFloat(), maxHeight = in.getFloat();
                String text = readRequiredString(in);
                String imageSource = readRequiredString(in);
                if (in.remaining() < 4) return null;
                int styleLength = in.getInt();
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return null;
                byte[] typedStyle = new byte[styleLength];
                in.get(typedStyle);
                measured.add(measureNativeControl(id, kind, text, imageSource, typedStyle, minWidth, maxWidth, minHeight, maxHeight));
            }
            if (in.hasRemaining()) return null;
            int capacity = 6;
            for (NativeMeasurement item : measured) capacity += 20 + item.error.getBytes(StandardCharsets.UTF_8).length;
            ByteBuffer out = ByteBuffer.allocate(capacity).order(ByteOrder.LITTLE_ENDIAN);
            out.putShort((short) 1).putInt(measured.size());
            for (NativeMeasurement item : measured) {
                byte[] error = item.error.getBytes(StandardCharsets.UTF_8);
                out.putLong(item.id).putFloat(item.width).putFloat(item.height).putInt(error.length).put(error);
            }
            return out.array();
        } catch (Throwable error) {
            android.util.Log.e("GoNative", "Native measurement batch failed", error);
            return null;
        }
    }

    private NativeMeasurement measureNativeControl(long id, int kind, String text, String imageSource, byte[] typedStyle,
                                                    float minWidth, float maxWidth, float minHeight, float maxHeight) {
        try {
            View view = makeView(kind, false);
            if (view instanceof TextView) {
                TextView label = (TextView) view;
                label.setText(text);
                label.setIncludeFontPadding(false);
                if (view instanceof Button) {
                    Button button = (Button) view;
                    button.setAllCaps(false);
                    button.setMinWidth(0);
                    button.setMinimumWidth(0);
                    button.setMinHeight(dp(44));
                    button.setMinimumHeight(dp(44));
                    button.setPadding(dp(16), 0, dp(16), 0);
                }
                if (view instanceof EditText) {
                    EditText field = (EditText) view;
                    field.setSingleLine(true);
                    field.setMinWidth(dp(240));
                    field.setMinimumWidth(dp(240));
                    field.setMinHeight(dp(44));
                    field.setMinimumHeight(dp(44));
                    field.setPadding(dp(12), dp(8), dp(12), dp(8));
                }
            }
            if (view instanceof ImageView && imageSource != null && !imageSource.isEmpty()) {
                int resource = getResources().getIdentifier(imageSource, "drawable", getPackageName());
                if (resource == 0) resource = getResources().getIdentifier(imageSource, "mipmap", getPackageName());
                if (resource != 0) ((ImageView) view).setImageResource(resource);
            }
            applyTypedStyle(view, typedStyle);
            int widthSpec = nativeMeasureSpec(maxWidth);
            int heightSpec = nativeMeasureSpec(maxHeight);
            view.measure(widthSpec, heightSpec);
            float density = getResources().getDisplayMetrics().density;
            float width = Math.max(finiteNonNegative(minWidth), view.getMeasuredWidth() / density);
            float height = Math.max(finiteNonNegative(minHeight), view.getMeasuredHeight() / density);
            if (isFinite(maxWidth) && maxWidth >= 0) width = Math.min(width, maxWidth);
            if (isFinite(maxHeight) && maxHeight >= 0) height = Math.min(height, maxHeight);
            return new NativeMeasurement(id, width, height, "");
        } catch (Throwable error) {
            return new NativeMeasurement(id, 0, 0, error.getClass().getSimpleName() + ": " + String.valueOf(error.getMessage()));
        }
    }

    private int nativeMeasureSpec(float maximum) {
        if (!isFinite(maximum) || maximum <= 0) return View.MeasureSpec.makeMeasureSpec(0, View.MeasureSpec.UNSPECIFIED);
        return View.MeasureSpec.makeMeasureSpec(dp(maximum), View.MeasureSpec.AT_MOST);
    }

    private static boolean isFinite(float value) { return !Float.isNaN(value) && !Float.isInfinite(value); }
    private static float finiteNonNegative(float value) { return isFinite(value) ? Math.max(0, value) : 0; }

    private static final class NativeMeasurement {
        final long id; final float width, height; final String error;
        NativeMeasurement(long id, float width, float height, String error) { this.id = id; this.width = width; this.height = height; this.error = error == null ? "" : error; }
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
        try {
            long started = System.nanoTime();
            ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
            if (in.remaining() < 14) return;
            if (Short.toUnsignedInt(in.getShort()) != 9) return;
            int count = in.getInt();
            long sequence = in.getLong();
            for (int operation = 0; operation < count && in.hasRemaining(); operation++) {
                int mutation = Byte.toUnsignedInt(in.get());
                int kind = Byte.toUnsignedInt(in.get());
                long nodeID = in.getLong();
                long parentID = in.getLong();
                int index = in.getInt();
                int fromIndex = in.getInt();
                float width = in.getFloat();
                float height = in.getFloat();
                float padding = in.getFloat();
                float gap = in.getFloat();
                int alignment = Byte.toUnsignedInt(in.get());
                boolean bold = in.get() != 0;
                float fontSize = in.getFloat();
                long handler = in.getLong();
                long changeHandler = in.getLong();
                long toggleHandler = in.getLong();
                boolean checked = in.get() != 0;
                float progress = in.getFloat();
                String text = readString(in);
                String accessibility = readString(in);
                String hint = readString(in);
                int role = Byte.toUnsignedInt(in.get());
                boolean focused = in.get() != 0;
                boolean scalesText = in.get() != 0;
                String imageSource = readString(in);
                int imageMode = Byte.toUnsignedInt(in.get());
                boolean horizontal = in.get() != 0;
                int interactionLength = in.remaining() >= 4 ? in.getInt() : 0;
                byte[] interactions = new byte[Math.max(0, Math.min(interactionLength, in.remaining()))];
                if (interactions.length > 0) in.get(interactions);
                int styleLength = in.remaining() >= 4 ? in.getInt() : -1;
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return;
                byte[] typedStyle = new byte[styleLength];
                in.get(typedStyle);
                if (in.remaining() < 17) return;
                boolean hasFrame = in.get() != 0;
                float frameX = in.getFloat(), frameY = in.getFloat(), frameWidth = in.getFloat(), frameHeight = in.getFloat();
                View view = views.get(nodeID);

                if (mutation == CREATE) {
                    view = makeView(kind, horizontal);
                    view.setTag(nodeID);
                    final long focusNodeID = nodeID;
                    view.setOnFocusChangeListener(new View.OnFocusChangeListener() {
                        @Override public void onFocusChange(View changed, boolean hasFocus) { nativeDispatchFocus(focusNodeID, hasFocus); }
                    });
                    views.put(nodeID, view);
                    if (rootNodeID == 0) rootNodeID = nodeID;
                    style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                    applyTypedStyle(view, typedStyle);
                    applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight);
                    applyInteractions(nodeID, view, interactions);
                    if (views.size() == 1) {
                        view.setLayoutParams(new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
                        view.setBackgroundColor(android.graphics.Color.WHITE);
                        setContentView(view);
                    }
                } else if (mutation == UPDATE) {
                    if (view != null) {
                        style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                        applyTypedStyle(view, typedStyle);
                        applyComputedFrame(nodeID, view, hasFrame, frameX, frameY, frameWidth, frameHeight);
                        applyInteractions(nodeID, view, interactions);
                    }
                } else if (mutation == INSERT) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == REMOVE) {
                    detach(view);
                } else if (mutation == MOVE) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        if (parent instanceof ScrollView || parent instanceof HorizontalScrollView) {
                            parent.removeAllViews();
                            ViewGroup.LayoutParams lp = new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
                            parent.addView(view, lp);
                        } else {
                            parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                        }
                    }
                } else if (mutation == DELETE) {
                    detach(view);
                    if (view instanceof EditText) textWatchers.remove((EditText) view);
                    GestureBinding binding = gestureBindings.get(nodeID);
                    if (binding != null) binding.dispose();
                    gestureBindings.remove(nodeID);
                    if (view != null) view.animate().cancel();
                    views.remove(nodeID);
                    if (nodeID == rootNodeID) rootNodeID = 0;
                }
                // Retained in the cross-platform protocol for renderer diagnostics.
                if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
            }
            nativeReportBatchApplied(sequence, System.nanoTime() - started);
        } catch (Throwable t) {
            android.util.Log.e("GoNative", "Error applying mutation batch", t);
        }
    }

    private void applyComputedFrame(long nodeID, View view, boolean hasFrame, float x, float y, float width, float height) {
        if (!hasFrame || view == null || nodeID == rootNodeID) return;
        if (!isFinite(x) || !isFinite(y) || !isFinite(width) || !isFinite(height) || width < 0 || height < 0) return;
        int measuredWidth = dp(Math.min(width, 1000000f));
        int measuredHeight = dp(Math.min(height, 1000000f));
        ViewGroup.LayoutParams params = view.getLayoutParams();
        if (params == null) params = new ViewGroup.LayoutParams(measuredWidth, measuredHeight);
        params.width = measuredWidth;
        params.height = measuredHeight;
        view.setLayoutParams(params);
        final float targetX = Math.max(-1000000f, Math.min(1000000f, x));
        final float targetY = Math.max(-1000000f, Math.min(1000000f, y));
        view.post(new Runnable() {
            @Override public void run() { view.setX(dp(targetX)); view.setY(dp(targetY)); }
        });
    }

    private View makeView(int kind, boolean horizontal) {
        if (kind == TEXT) {
            TextView tv = new TextView(this);
            tv.setTextColor(android.graphics.Color.BLACK);
            return tv;
        }
        if (kind == BUTTON) return new Button(this);
        if (kind == TEXT_INPUT) return new EditText(this);
        if (kind == SWITCH) return new Switch(this);
        if (kind == PROGRESS_INDICATOR) { ProgressBar bar = new ProgressBar(this, null, android.R.attr.progressBarStyleHorizontal); bar.setMax(10000); return bar; }
        if (kind == IMAGE) return new ImageView(this);
        if (kind == SCROLL_VIEW) {
            if (horizontal) {
                HorizontalScrollView hsv = new HorizontalScrollView(this);
                hsv.setFillViewport(true);
                return hsv;
            } else {
                ScrollView sv = new ScrollView(this);
                sv.setFillViewport(true);
                return sv;
            }
        }
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    private void applyTypedStyle(View view, byte[] payload) {
        if (view == null || payload == null || payload.length < 187) return;
        ByteBuffer style = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
        if (Short.toUnsignedInt(style.getShort(0)) != 1) return;
        int portable = 2, ios = portable + typedStyleSize(style, portable, payload.length), androidStyle = ios + typedStyleSize(style, ios, payload.length);
        if (androidStyle <= ios || androidStyle >= payload.length || typedStyleSize(style, androidStyle, payload.length) == 0) return;
        int appearanceBase = hasTypedValues(payload, androidStyle + 112, 69) ? androidStyle + 112 : portable + 112;
        int androidFontLength = style.getInt(androidStyle + 181);
        int textBase = hasTypedValues(payload, androidStyle + 181, 22 + Math.max(0, androidFontLength)) ? androidStyle : portable;
        int fontLength = style.getInt(textBase + 181), interactionBase = hasTypedValues(payload, androidStyle + 203 + Math.max(0, androidFontLength), 17) ? androidStyle : portable;
        int background = rgba(style, appearanceBase), foreground = rgba(style, appearanceBase + 4);
        float borderWidth = style.getFloat(appearanceBase + 8), cornerRadius = style.getFloat(appearanceBase + 16), opacity = style.getFloat(appearanceBase + 44);
        int borderColor = rgba(style, appearanceBase + 12), visibility = Byte.toUnsignedInt(style.get(appearanceBase + 68));
        int disabledOffset = interactionBase + 203 + style.getInt(interactionBase + 181);
        if (fontLength < 0 || disabledOffset >= payload.length) return;
        int fontOffset = textBase + 185 + fontLength;
        float fontSize = style.getFloat(fontOffset), lineHeight = style.getFloat(fontOffset + 6), letterSpacing = style.getFloat(fontOffset + 10);
        int fontWeight = Short.toUnsignedInt(style.getShort(fontOffset + 4));
        if (view instanceof TextView) {
            TextView text = (TextView) view;
            String family = new String(payload, textBase + 185, fontLength, java.nio.charset.StandardCharsets.UTF_8);
            Typeface face = family.isEmpty() ? Typeface.DEFAULT : Typeface.create(family, Typeface.NORMAL);
            text.setTypeface(face, fontWeight >= 600 ? Typeface.BOLD : Typeface.NORMAL);
            if (fontSize > 0) text.setTextSize(TypedValue.COMPLEX_UNIT_DIP, fontSize);
            if (lineHeight > 0 && android.os.Build.VERSION.SDK_INT >= 28) text.setLineHeight(dp(lineHeight));
            if (letterSpacing != 0 && fontSize > 0) text.setLetterSpacing(letterSpacing / fontSize);
        }
        if (android.graphics.Color.alpha(background) > 0 || borderWidth > 0) {
            android.graphics.drawable.GradientDrawable drawable = new android.graphics.drawable.GradientDrawable();
            drawable.setColor(background);
            if (cornerRadius > 0) drawable.setCornerRadius(dp(cornerRadius));
            if (borderWidth > 0) drawable.setStroke(dp(borderWidth), borderColor);
            view.setBackground(drawable);
        }
        if (view instanceof TextView && android.graphics.Color.alpha(foreground) > 0) ((TextView) view).setTextColor(foreground);
        if (opacity > 0) view.setAlpha(Math.min(1f, opacity));
        float translateX = style.getFloat(appearanceBase + 48), translateY = style.getFloat(appearanceBase + 52), scaleX = style.getFloat(appearanceBase + 56), scaleY = style.getFloat(appearanceBase + 60), rotation = style.getFloat(appearanceBase + 64);
        view.setTranslationX(dp(translateX)); view.setTranslationY(dp(translateY));
        view.setScaleX(scaleX == 0 ? 1 : scaleX); view.setScaleY(scaleY == 0 ? 1 : scaleY); view.setRotation(rotation);
        float shadowBlur = style.getFloat(appearanceBase + 32), shadowOpacity = style.getFloat(appearanceBase + 40);
        if (shadowBlur > 0 && shadowOpacity > 0) view.setElevation(dp(shadowBlur));
        view.setVisibility(visibility == 2 ? View.GONE : visibility == 1 ? View.INVISIBLE : View.VISIBLE);
        view.setEnabled(style.get(disabledOffset) == 0);
    }

    private int rgba(ByteBuffer style, int offset) {
        return android.graphics.Color.argb(Byte.toUnsignedInt(style.get(offset + 3)), Byte.toUnsignedInt(style.get(offset)), Byte.toUnsignedInt(style.get(offset + 1)), Byte.toUnsignedInt(style.get(offset + 2)));
    }

    private int typedStyleSize(ByteBuffer style, int base, int limit) {
        if (base < 0 || base + 185 > limit) return 0;
        int fontLength = style.getInt(base + 181);
        int size = 220 + fontLength;
        return fontLength < 0 || base + size > limit ? 0 : size;
    }

    private boolean hasTypedValues(byte[] payload, int offset, int length) {
        if (offset < 0 || length < 0 || offset + length > payload.length) return false;
        for (int i = offset; i < offset + length; i++) if (payload[i] != 0) return true;
        return false;
    }

    private void style(View view, int kind, String text, float width, float height, float padding,
                       float gap, int alignment, float fontSize, boolean bold, long handler, long changeHandler, long toggleHandler, boolean checked, float progress,
                       String accessibility, String hint, int role, boolean focused, boolean scalesText, String imageSource, int imageMode) {
        if (kind == TEXT && view instanceof TextView) {
            TextView textView = (TextView) view;
            textView.setText(text);
            if (fontSize > 0) textView.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            textView.setIncludeFontPadding(false);
            textView.setTextColor(android.graphics.Color.parseColor("#111111"));
        }
        if (kind == BUTTON && view instanceof Button) {
            Button btn = (Button) view;
            btn.setText(text);
            if (fontSize > 0) btn.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            btn.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            btn.setAllCaps(false);
            btn.setIncludeFontPadding(false);
            btn.setGravity(Gravity.CENTER);
            btn.setMinHeight(0);
            btn.setMinimumHeight(0);
            btn.setElevation(0);
            android.graphics.drawable.GradientDrawable btnBg = new android.graphics.drawable.GradientDrawable();
            btnBg.setColor(android.graphics.Color.parseColor("#007AFF"));
            btnBg.setCornerRadius(dp(8));
            btn.setBackground(btnBg);
            btn.setTextColor(android.graphics.Color.WHITE);
            btn.setPadding(dp(16), 0, dp(16), 0);
            final long eventHandler = handler;
            if (eventHandler != 0) {
                view.setOnClickListener(new View.OnClickListener() {
                    @Override public void onClick(View clicked) { nativeDispatchEvent(eventHandler); }
                });
            } else {
                view.setOnClickListener(null);
            }
        }
        if (view instanceof EditText) {
            EditText field = (EditText) view;
            field.setSingleLine(true);
            field.setGravity(Gravity.CENTER_VERTICAL);
            field.setIncludeFontPadding(false);
            field.setMinHeight(dp(44));
            if (hint != null && !hint.isEmpty()) field.setHint(hint);
            field.setTextColor(android.graphics.Color.BLACK);
            field.setHintTextColor(android.graphics.Color.parseColor("#8E8E93"));
            android.graphics.drawable.GradientDrawable fieldBg = new android.graphics.drawable.GradientDrawable();
            fieldBg.setColor(android.graphics.Color.parseColor("#FAFAFC"));
            fieldBg.setCornerRadius(dp(8));
            fieldBg.setStroke(dp(1), android.graphics.Color.parseColor("#D1D1D6"));
            field.setBackground(fieldBg);
            int padX = dp(12), padY = dp(8);
            field.setPadding(padX, padY, padX, padY);
            TextWatcher existing = textWatchers.get(field);
            if (existing != null) field.removeTextChangedListener(existing);
            if (text != null && !field.getText().toString().equals(text)) { field.setText(text); field.setSelection(field.length()); }
            field.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize > 0 ? fontSize : 16);
            field.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            final long eventHandler = changeHandler;
            TextWatcher watcher = new TextWatcher() {
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                public void onTextChanged(CharSequence s, int start, int before, int count) { if (eventHandler != 0) nativeDispatchValueEvent(eventHandler, s.toString()); }
                public void afterTextChanged(Editable s) {}
            };
            field.addTextChangedListener(watcher); textWatchers.put(field, watcher);
        }
        if (view instanceof Switch) {
            Switch toggle = (Switch) view;
            toggle.setOnCheckedChangeListener(null); toggle.setChecked(checked);
            final long eventHandler = toggleHandler;
            if (eventHandler != 0) {
                toggle.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
                    @Override public void onCheckedChanged(CompoundButton button, boolean value) { nativeDispatchBoolEvent(eventHandler, value); }
                });
            }
        }
        if (view instanceof ProgressBar) { ((ProgressBar) view).setProgress(Math.round(progress * 10000)); }
        if (view instanceof ImageView) {
            ImageView image = (ImageView) view;
            if (imageSource != null && !imageSource.isEmpty()) {
                int resource = getResources().getIdentifier(imageSource, "drawable", getPackageName());
                if (resource == 0) resource = getResources().getIdentifier(imageSource, "mipmap", getPackageName());
                if (resource != 0) {
                    image.setImageResource(resource);
                } else if ("app_logo".equals(imageSource)) {
                    android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setCornerRadius(dp(16));
                    image.setBackground(gd);
                    image.setImageResource(android.R.drawable.ic_lock_lock);
                    image.setColorFilter(android.graphics.Color.WHITE);
                    int pad = dp(12);
                    image.setPadding(pad, pad, pad, pad);
                } else if ("avatar".equals(imageSource)) {
                    android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setShape(android.graphics.drawable.GradientDrawable.OVAL);
                    image.setBackground(gd);
                    image.setImageResource(android.R.drawable.ic_menu_myplaces);
                    image.setColorFilter(android.graphics.Color.WHITE);
                    int pad = dp(10);
                    image.setPadding(pad, pad, pad, pad);
                } else {
                    image.setImageDrawable(null);
                }
            } else {
                image.setImageDrawable(null);
            }
            image.setScaleType(imageMode == 1 ? ImageView.ScaleType.CENTER_CROP : imageMode == 2 ? ImageView.ScaleType.CENTER : ImageView.ScaleType.FIT_CENTER);
        }
        int paddingPx = dp(padding);
        if (view instanceof EditText) {
            view.setPadding(dp(12) + paddingPx, dp(8) + paddingPx, dp(12) + paddingPx, dp(8) + paddingPx);
        } else {
            view.setPadding(paddingPx, paddingPx, paddingPx, paddingPx);
        }
        view.setContentDescription(accessibility.isEmpty() ? text : accessibility);
        view.setAccessibilityDelegate(new View.AccessibilityDelegate() {
            @Override public void onInitializeAccessibilityNodeInfo(View host, AccessibilityNodeInfo info) {
                super.onInitializeAccessibilityNodeInfo(host, info);
                if (role == 1 || role == 3) info.setClassName(TextView.class.getName());
                else if (role == 2) info.setClassName(Button.class.getName());
                else if (role == 4) info.setClassName("android.widget.ImageView");
                if (android.os.Build.VERSION.SDK_INT >= 26 && !hint.isEmpty()) info.setHintText(hint);
                if (android.os.Build.VERSION.SDK_INT >= 28 && role == 3) info.setHeading(true);
            }
        });
        if (focused) {
            view.requestFocus();
            view.sendAccessibilityEvent(AccessibilityEvent.TYPE_VIEW_FOCUSED);
        } else if (view.hasFocus()) {
            view.clearFocus();
        }
        ViewGroup.LayoutParams current = view.getLayoutParams();
        int requestedWidth = width > 0 ? dp(width) : ViewGroup.LayoutParams.WRAP_CONTENT;
        int requestedHeight = height > 0 ? dp(height) : ViewGroup.LayoutParams.WRAP_CONTENT;
        if (current == null) current = new LinearLayout.LayoutParams(requestedWidth, requestedHeight);
        current.width = requestedWidth;
        current.height = requestedHeight;
        view.setLayoutParams(current);
        if (view instanceof LinearLayout) {
            LinearLayout container = (LinearLayout) view;
            container.setGravity(alignment == 1 ? Gravity.CENTER : alignment == 2 ? Gravity.END : Gravity.START);
            container.setShowDividers(gap > 0 ? LinearLayout.SHOW_DIVIDER_MIDDLE : LinearLayout.SHOW_DIVIDER_NONE);
            if (gap > 0) container.setDividerDrawable(new GapDrawable(dp(gap), container.getOrientation()));
        }
    }

    private int dp(float value) { return Math.round(value * getResources().getDisplayMetrics().density); }

    private void applyInteractions(long nodeID, View view, byte[] payload) {
        if (view == null) return;
        GestureBinding old = gestureBindings.get(nodeID);
        if (old != null) old.dispose();
        gestureBindings.remove(nodeID);
        if (payload == null || payload.length < 4) return;
        ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
        ArrayList<GestureSpec> gestures = new ArrayList<>();
        int gestureCount = in.getInt();
        for (int i = 0; i < gestureCount && in.remaining() >= 26; i++) {
            GestureSpec spec = new GestureSpec();
            spec.kind = Byte.toUnsignedInt(in.get());
            spec.direction = Byte.toUnsignedInt(in.get());
            spec.minimumPressNanos = in.getLong();
            spec.minimumTravel = in.getFloat();
            spec.handler = in.getLong();
            gestures.add(spec);
        }
        if (!gestures.isEmpty()) {
            GestureBinding binding = new GestureBinding(view, gestures);
            gestureBindings.put(nodeID, binding);
            view.setOnTouchListener(binding);
        } else view.setOnTouchListener(null);

        if (in.remaining() < 4) return;
        int animationCount = in.getInt();
        ArrayList<Animator> animations = new ArrayList<>();
        for (int i = 0; i < animationCount && in.remaining() >= 42; i++) {
            int property = Byte.toUnsignedInt(in.get());
            long durationNanos = in.getLong();
            long delayNanos = in.getLong();
            int curve = Byte.toUnsignedInt(in.get());
            float damping = in.getFloat();
            float velocity = in.getFloat();
            boolean reduceMotionOK = in.get() != 0;
            float from = in.getFloat(), to = in.getFloat(), fromX = in.getFloat(), fromY = in.getFloat(), toX = in.getFloat(), toY = in.getFloat();
            Animator animator = makeAnimator(view, property, from, to, fromX, fromY, toX, toY);
            if (animator == null) continue;
            animator.setDuration(Math.max(0, durationNanos / 1000000L));
            animator.setStartDelay(Math.max(0, delayNanos / 1000000L));
            if (curve == 1) animator.setInterpolator(new AccelerateInterpolator());
            else if (curve == 2) animator.setInterpolator(new DecelerateInterpolator());
            else if (curve == 3) animator.setInterpolator(new LinearInterpolator());
            else if (curve == 4) animator.setInterpolator(new OvershootInterpolator(Math.max(.1f, (1f - damping) * 2f + Math.abs(velocity) * .1f)));
            else animator.setInterpolator(new AccelerateDecelerateInterpolator());
            if (!reduceMotionOK && android.os.Build.VERSION.SDK_INT >= 26 && !ValueAnimator.areAnimatorsEnabled()) animator.setDuration(0);
            animations.add(animator);
        }
        if (!animations.isEmpty()) { AnimatorSet set = new AnimatorSet(); set.playSequentially(animations); set.start(); }
    }

    private Animator makeAnimator(final View view, int property, float from, float to, float fromX, float fromY, float toX, float toY) {
        if (property == 1) return ObjectAnimator.ofFloat(view, View.ALPHA, from, to);
        if (property == 2) { AnimatorSet set = new AnimatorSet(); set.playTogether(ObjectAnimator.ofFloat(view, View.SCALE_X, from, to), ObjectAnimator.ofFloat(view, View.SCALE_Y, from, to)); return set; }
        if (property == 3) { AnimatorSet set = new AnimatorSet(); set.playTogether(ObjectAnimator.ofFloat(view, View.TRANSLATION_X, dp(fromX), dp(toX)), ObjectAnimator.ofFloat(view, View.TRANSLATION_Y, dp(fromY), dp(toY))); return set; }
        if (property == 4) {
            ValueAnimator animator = ValueAnimator.ofFloat(0, 1);
            animator.addUpdateListener(new ValueAnimator.AnimatorUpdateListener() {
                @Override public void onAnimationUpdate(ValueAnimator animation) { view.requestLayout(); }
            });
            return animator;
        }
        return null;
    }

    private static final class GestureSpec { int kind, direction; long minimumPressNanos, handler; float minimumTravel; }

    private final class GestureBinding implements View.OnTouchListener {
        final View view; final List<GestureSpec> specs; final ArrayList<Runnable> pending = new ArrayList<>();
        VelocityTracker velocity; float downX, downY; long downTime; boolean moved;
        GestureBinding(View view, List<GestureSpec> specs) { this.view = view; this.specs = specs; }
        void dispose() { for (Runnable r : pending) view.removeCallbacks(r); pending.clear(); if (velocity != null) velocity.recycle(); velocity = null; view.setOnTouchListener(null); }
        void emit(GestureSpec s, float x, float y, float vx, float vy) { if (s.handler != 0) nativeDispatchGestureEvent(s.handler, x / getResources().getDisplayMetrics().density, y / getResources().getDisplayMetrics().density, vx / getResources().getDisplayMetrics().density, vy / getResources().getDisplayMetrics().density); }
        @Override public boolean onTouch(View ignored, MotionEvent event) {
            if (event.getActionMasked() == MotionEvent.ACTION_DOWN) {
                downX=event.getX(); downY=event.getY(); downTime=System.nanoTime(); moved=false;
                velocity=VelocityTracker.obtain(); velocity.addMovement(event);
                for (final GestureSpec s : specs) if (s.kind == 2) {
                    Runnable r = new Runnable() {
                        @Override public void run() { if (!moved) emit(s, 0, 0, 0, 0); }
                    };
                    pending.add(r);
                    view.postDelayed(r, Math.max(0, s.minimumPressNanos / 1000000L));
                }
                return true;
            }
            if (velocity != null) velocity.addMovement(event);
            float dx=event.getX()-downX, dy=event.getY()-downY;
            if (event.getActionMasked() == MotionEvent.ACTION_MOVE) {
                for (GestureSpec s:specs) if (Math.hypot(dx,dy) >= dp(s.minimumTravel)) { moved=true; if(s.kind==4) { velocity.computeCurrentVelocity(1000); emit(s,dx,dy,velocity.getXVelocity(),velocity.getYVelocity()); } }
                if(moved) { for(Runnable r:pending)view.removeCallbacks(r); pending.clear(); }
            } else if (event.getActionMasked() == MotionEvent.ACTION_UP) {
                for(Runnable r:pending)view.removeCallbacks(r); pending.clear(); velocity.computeCurrentVelocity(1000); float vx=velocity.getXVelocity(),vy=velocity.getYVelocity();
                for(GestureSpec s:specs) {
                    float distance=(float)Math.hypot(dx,dy);
                    if(s.kind==1 && distance < dp(Math.max(8,s.minimumTravel))) emit(s,dx,dy,vx,vy);
                    if(s.kind==3 && distance >= dp(s.minimumTravel) && directionMatches(s.direction,dx,dy)) emit(s,dx,dy,vx,vy);
                    if(s.kind==4) emit(s,dx,dy,vx,vy);
                }
                if (view instanceof Button && !moved) view.performClick();
                velocity.recycle(); velocity=null;
            } else if(event.getActionMasked()==MotionEvent.ACTION_CANCEL) disposePending();
            return true;
        }
        void disposePending(){for(Runnable r:pending)view.removeCallbacks(r);pending.clear();if(velocity!=null){velocity.recycle();velocity=null;}}
        boolean directionMatches(int direction,float dx,float dy){ if(direction==0)return true;if(direction==1)return dy<0&&Math.abs(dy)>=Math.abs(dx);if(direction==2)return dy>0&&Math.abs(dy)>=Math.abs(dx);boolean rtl=view.getLayoutDirection()==View.LAYOUT_DIRECTION_RTL;if(direction==3)return rtl?dx>0:dx<0;if(direction==4)return rtl?dx<0:dx>0;return false; }
    }

    private static void detach(View view) {
        if (view != null && view.getParent() instanceof ViewGroup) ((ViewGroup) view.getParent()).removeView(view);
    }

    private static String readString(ByteBuffer in) {
        if (in.remaining() < 4) return "";
        int length = in.getInt();
        if (length <= 0 || in.remaining() < length) return "";
        byte[] bytes = new byte[length];
        in.get(bytes);
        return new String(bytes, StandardCharsets.UTF_8);
    }

    private static String readRequiredString(ByteBuffer in) {
        if (in.remaining() < 4) throw new IllegalArgumentException("missing string length");
        int length = in.getInt();
        if (length < 0 || length > 1048576 || length > in.remaining()) throw new IllegalArgumentException("invalid string length");
        byte[] bytes = new byte[length];
        in.get(bytes);
        return new String(bytes, StandardCharsets.UTF_8);
    }
}
