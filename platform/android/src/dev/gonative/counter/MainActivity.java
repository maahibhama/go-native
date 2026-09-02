package dev.gonative.counter;

import android.app.Activity;
import android.graphics.Typeface;
import android.os.Bundle;
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

    private native void nativeStart();
    private native void nativeDispatchEvent(long handler);
    private native void nativeDispatchValueEvent(long handler, String value);
    private native void nativeDispatchBoolEvent(long handler, boolean value);
    private native void nativeDispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY);
    private native void nativeStop();
    private native void nativeReportBatchApplied(long sequence, long nativeNanos);

    @Override protected void onCreate(Bundle state) {
        super.onCreate(state);
        nativeStart();
    }

    @Override protected void onDestroy() {
        nativeStop();
        for (int i = 0; i < gestureBindings.size(); i++) gestureBindings.valueAt(i).dispose();
        gestureBindings.clear();
        for (int i = 0; i < views.size(); i++) views.valueAt(i).animate().cancel();
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
        if (payload == null) return;
        try {
            long started = System.nanoTime();
            ByteBuffer in = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
            if (in.remaining() < 14) return;
            if (Short.toUnsignedInt(in.getShort()) != 7) return;
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
                View view = views.get(nodeID);

                if (mutation == CREATE) {
                    view = makeView(kind, horizontal);
                    view.setTag(nodeID);
                    views.put(nodeID, view);
                    style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                    applyInteractions(nodeID, view, interactions);
                    if (views.size() == 1) setContentView(view);
                } else if (mutation == UPDATE) {
                    if (view != null) {
                        style(view, kind, text, width, height, padding, gap, alignment, fontSize, bold, handler, changeHandler, toggleHandler, checked, progress, accessibility, hint, role, focused, scalesText, imageSource, imageMode);
                        applyInteractions(nodeID, view, interactions);
                    }
                } else if (mutation == INSERT) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                    }
                } else if (mutation == REMOVE) {
                    detach(view);
                } else if (mutation == MOVE) {
                    View parentView = views.get(parentID);
                    if (parentView instanceof ViewGroup && view != null) {
                        ViewGroup parent = (ViewGroup) parentView;
                        detach(view);
                        parent.addView(view, Math.min(Math.max(index, 0), parent.getChildCount()));
                    }
                } else if (mutation == DELETE) {
                    detach(view);
                    GestureBinding binding = gestureBindings.get(nodeID);
                    if (binding != null) binding.dispose();
                    gestureBindings.remove(nodeID);
                    if (view != null) view.animate().cancel();
                    views.remove(nodeID);
                }
                // Retained in the cross-platform protocol for renderer diagnostics.
                if (fromIndex == Integer.MIN_VALUE) throw new AssertionError();
            }
            nativeReportBatchApplied(sequence, System.nanoTime() - started);
        } catch (Throwable t) {
            android.util.Log.e("GoNative", "Error applying mutation batch", t);
        }
    }

    private View makeView(int kind, boolean horizontal) {
        if (kind == TEXT) return new TextView(this);
        if (kind == BUTTON) return new Button(this);
        if (kind == TEXT_INPUT) return new EditText(this);
        if (kind == SWITCH) return new Switch(this);
        if (kind == PROGRESS_INDICATOR) { ProgressBar bar = new ProgressBar(this, null, android.R.attr.progressBarStyleHorizontal); bar.setMax(10000); return bar; }
        if (kind == IMAGE) return new ImageView(this);
        if (kind == SCROLL_VIEW) return horizontal ? new HorizontalScrollView(this) : new ScrollView(this);
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    private void style(View view, int kind, String text, float width, float height, float padding,
                       float gap, int alignment, float fontSize, boolean bold, long handler, long changeHandler, long toggleHandler, boolean checked, float progress,
                       String accessibility, String hint, int role, boolean focused, boolean scalesText, String imageSource, int imageMode) {
        if (kind == TEXT && view instanceof TextView) {
            TextView textView = (TextView) view;
            textView.setText(text);
            if (fontSize > 0) textView.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
        }
        if (kind == BUTTON && view instanceof Button) {
            Button btn = (Button) view;
            btn.setText(text);
            if (fontSize > 0) btn.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            btn.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
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
            Object existing = field.getTag(android.R.id.custom);
            if (existing instanceof TextWatcher) field.removeTextChangedListener((TextWatcher) existing);
            if (text != null && !field.getText().toString().equals(text)) { field.setText(text); field.setSelection(field.length()); }
            if (fontSize > 0) field.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            field.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            final long eventHandler = changeHandler;
            TextWatcher watcher = new TextWatcher() {
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                public void onTextChanged(CharSequence s, int start, int before, int count) { if (eventHandler != 0) nativeDispatchValueEvent(eventHandler, s.toString()); }
                public void afterTextChanged(Editable s) {}
            };
            field.addTextChangedListener(watcher); field.setTag(android.R.id.custom, watcher);
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
                if (resource != 0) image.setImageResource(resource);
                else image.setImageDrawable(null);
            } else {
                image.setImageDrawable(null);
            }
            image.setScaleType(imageMode == 1 ? ImageView.ScaleType.CENTER_CROP : imageMode == 2 ? ImageView.ScaleType.CENTER : ImageView.ScaleType.FIT_CENTER);
        }
        int paddingPx = dp(padding);
        view.setPadding(paddingPx, paddingPx, paddingPx, paddingPx);
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
}
