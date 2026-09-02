package dev.gonative.counter;

import android.graphics.Canvas;
import android.graphics.ColorFilter;
import android.graphics.PixelFormat;
import android.graphics.drawable.Drawable;
import android.widget.LinearLayout;

@SuppressWarnings("deprecation")
final class GapDrawable extends Drawable {
    private final int size;
    private final int orientation;

    GapDrawable(int size, int orientation) {
        this.size = size;
        this.orientation = orientation;
    }

    @Override public void draw(Canvas canvas) {}
    @Override public void setAlpha(int alpha) {}
    @Override public void setColorFilter(ColorFilter filter) {}
    @Override public int getOpacity() { return PixelFormat.TRANSPARENT; }
    @Override public int getIntrinsicWidth() { return orientation == LinearLayout.HORIZONTAL ? size : 0; }
    @Override public int getIntrinsicHeight() { return orientation == LinearLayout.VERTICAL ? size : 0; }
}
