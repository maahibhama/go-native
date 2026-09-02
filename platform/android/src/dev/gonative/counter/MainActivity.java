package dev.gonative.counter;

import android.app.Activity;
import android.graphics.Typeface;
import android.os.Bundle;
import android.util.LongSparseArray;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

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

    static { System.loadLibrary("gonative"); }

    private final LongSparseArray<View> views = new LongSparseArray<>();

    private native void nativeStart();
    private native void nativeDispatchEvent(long handler);
    private native void nativeStop();
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        nativeStart();
    }

    @Override protected void onDestroy() {
        nativeStop();
        views.clear();
        super.onDestroy();
    }

    @SuppressWarnings("unused")
    public void applyMutationBatch(byte[] payload) {
        final byte[] owned = payload.clone();
        runOnUiThread(new Runnable() {
            @Override public void run() { applyOnUiThread(owned); }
        });
    }

    private void applyOnUiThread(byte[] payload) {
        long started = System.nanoTime();
        ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
        if (Short.toUnsignedInt(in.getShort()) != 2) return;
        int count = in.getInt();
        long sequence = in.getLong();
        for (int operation = 0; operation < count; operation++) {
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
            String text = readString(in);
            String accessibility = readString(in);
            View view = views.get(nodeID);

            if (mutation == CREATE) {
                view = makeView(kind);
                view.setTag(nodeID);
                views.put(nodeID, view);
                style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, accessibility);
                if (views.size() == 1) setContentView(view);
            } else if (mutation == UPDATE) {
                style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, accessibility);
            } else if (mutation == INSERT) {
                ViewGroup parent = (ViewGroup) views.get(parentID);
                detach(view);
                parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
            } else if (mutation == REMOVE) {
                detach(view);
            } else if (mutation == MOVE) {
                ViewGroup parent = (ViewGroup) views.get(parentID);
                detach(view);
                parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
            } else if (mutation == DELETE) {
                detach(view);
                views.remove(nodeID);
            }
            // Retained in the cross-platform protocol for renderer diagnostics.
            if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
        }
        nativeReportBatchApplied(sequence, System.nanoTime() - started);
    }

    private View makeView(int kind) {
        if (kind == TEXT) return new TextView(this);
        if (kind == BUTTON) return new Button(this);
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    private void style(View view, int kind, String text, float width, float height, float padding,
                       float gap, int alignment, float fontSize, boolean bold, long handler,
                       String accessibility) {
        if (view instanceof TextView) {
            TextView textView = (TextView) view;
            textView.setText(text);
            if (fontSize > 0) textView.setTextSize(fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
        }
        if (view instanceof Button) {
            final long eventHandler = handler;
            view.setOnClickListener(new View.OnClickListener() {
                @Override public void onClick(View clicked) { nativeDispatchEvent(eventHandler); }
            });
        }
        int paddingPx = dp(padding);
        view.setPadding(paddingPx, paddingPx, paddingPx, paddingPx);
        view.setContentDescription(accessibility.isEmpty() ? text : accessibility);
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

    private static void detach(View view) {
        if (view != null && view.getParent() instanceof ViewGroup) ((ViewGroup) view.getParent()).removeView(view);
    }

    private static String readString(ByteBuffer in) {
        int length = in.getInt();
        byte[] bytes = new byte[length];
        in.get(bytes);
        return new String(bytes, StandardCharsets.UTF_8);
    }
}
