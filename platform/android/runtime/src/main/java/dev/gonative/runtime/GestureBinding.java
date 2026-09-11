package dev.gonative.runtime;

import android.view.MotionEvent;
import android.view.VelocityTracker;
import android.view.View;
import android.widget.Button;

import java.util.ArrayList;
import java.util.List;

public final class GestureBinding implements View.OnTouchListener {
    public static final class GestureSpec {
        public int kind, direction;
        public long minimumPressNanos, handler;
        public float minimumTravel;
    }

    final View view;
    final List<GestureSpec> specs;
    final EventDispatcher dispatcher;
    final ArrayList<Runnable> pending = new ArrayList<>();
    VelocityTracker velocity;
    float downX, downY;
    long downTime;
    boolean moved;

    public GestureBinding(View view, List<GestureSpec> specs, EventDispatcher dispatcher) {
        this.view = view;
        this.specs = specs;
        this.dispatcher = dispatcher;
    }

    public void dispose() {
        for (Runnable r : pending) view.removeCallbacks(r);
        pending.clear();
        if (velocity != null) velocity.recycle();
        velocity = null;
        view.setOnTouchListener(null);
    }

    private int dp(float value) {
        return Math.round(value * view.getResources().getDisplayMetrics().density);
    }

    private void emit(GestureSpec s, float x, float y, float vx, float vy) {
        if (s.handler != 0 && dispatcher != null) {
            float density = view.getResources().getDisplayMetrics().density;
            dispatcher.dispatchGestureEvent(s.handler, x / density, y / density, vx / density, vy / density);
        }
    }

    @Override
    public boolean onTouch(View ignored, MotionEvent event) {
        if (event.getActionMasked() == MotionEvent.ACTION_DOWN) {
            downX = event.getX();
            downY = event.getY();
            downTime = System.nanoTime();
            moved = false;
            velocity = VelocityTracker.obtain();
            velocity.addMovement(event);
            for (final GestureSpec s : specs) {
                if (s.kind == 2) {
                    Runnable r = new Runnable() {
                        @Override public void run() { if (!moved) emit(s, 0, 0, 0, 0); }
                    };
                    pending.add(r);
                    view.postDelayed(r, Math.max(0, s.minimumPressNanos / 1000000L));
                }
            }
            return true;
        }
        if (velocity != null) velocity.addMovement(event);
        float dx = event.getX() - downX, dy = event.getY() - downY;
        if (event.getActionMasked() == MotionEvent.ACTION_MOVE) {
            for (GestureSpec s : specs) {
                if (Math.hypot(dx, dy) >= dp(s.minimumTravel)) {
                    moved = true;
                    if (s.kind == 4) {
                        velocity.computeCurrentVelocity(1000);
                        emit(s, dx, dy, velocity.getXVelocity(), velocity.getYVelocity());
                    }
                }
            }
            if (moved) {
                for (Runnable r : pending) view.removeCallbacks(r);
                pending.clear();
            }
        } else if (event.getActionMasked() == MotionEvent.ACTION_UP) {
            for (Runnable r : pending) view.removeCallbacks(r);
            pending.clear();
            velocity.computeCurrentVelocity(1000);
            float vx = velocity.getXVelocity(), vy = velocity.getYVelocity();
            for (GestureSpec s : specs) {
                float distance = (float) Math.hypot(dx, dy);
                if (s.kind == 1 && distance < dp(Math.max(8, s.minimumTravel))) emit(s, dx, dy, vx, vy);
                if (s.kind == 3 && distance >= dp(s.minimumTravel) && directionMatches(s.direction, dx, dy)) emit(s, dx, dy, vx, vy);
                if (s.kind == 4) emit(s, dx, dy, vx, vy);
            }
            if (view instanceof Button && !moved) view.performClick();
            velocity.recycle();
            velocity = null;
        } else if (event.getActionMasked() == MotionEvent.ACTION_CANCEL) {
            disposePending();
        }
        return true;
    }

    private void disposePending() {
        for (Runnable r : pending) view.removeCallbacks(r);
        pending.clear();
        if (velocity != null) {
            velocity.recycle();
            velocity = null;
        }
    }

    private boolean directionMatches(int direction, float dx, float dy) {
        if (direction == 0) return true;
        if (direction == 1) return dy < 0 && Math.abs(dy) >= Math.abs(dx);
        if (direction == 2) return dy > 0 && Math.abs(dy) >= Math.abs(dx);
        boolean rtl = view.getLayoutDirection() == View.LAYOUT_DIRECTION_RTL;
        if (direction == 3) return rtl ? dx > 0 : dx < 0;
        if (direction == 4) return rtl ? dx < 0 : dx > 0;
        return false;
    }
}
