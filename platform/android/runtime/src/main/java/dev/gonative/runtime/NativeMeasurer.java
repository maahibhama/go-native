package dev.gonative.runtime;

import android.content.Context;
import android.os.Looper;
import android.text.InputFilter;
import android.text.TextUtils;
import android.text.method.PasswordTransformationMethod;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;

public final class NativeMeasurer {
    private final Context context;
    private final ControlFactory factory;

    public NativeMeasurer(Context context, ControlFactory factory) {
        this.context = context;
        this.factory = factory;
    }

    public static final class NativeMeasurement {
        public final long id;
        public final float width, height;
        public final String error;

        public NativeMeasurement(long id, float width, float height, String error) {
            this.id = id;
            this.width = width;
            this.height = height;
            this.error = error == null ? "" : error;
        }
    }

    public byte[] measureNativeBatchOnUiThread(byte[] payload) {
        try {
            ProtocolReader in = new ProtocolReader(payload);
            if (in.remaining() < 6 || in.uint16() != 2) return null;
            int count = in.int32();
            if (count < 0 || count > 100000) return null;
            ArrayList<NativeMeasurement> measured = new ArrayList<>(count);
            for (int i = 0; i < count; i++) {
                if (in.remaining() < 25) return null;
                long id = in.int64();
                int kind = in.uint8();
                float minWidth = in.float32(), maxWidth = in.float32(), minHeight = in.float32(), maxHeight = in.float32();
                String text = in.requiredString();
                String imageSource = in.requiredString();
                if (in.remaining() < 13) return null;
                int textWrap = in.uint8(), textOverflow = in.uint8();
                long maxLinesValue = in.uint32();
                if (maxLinesValue > Integer.MAX_VALUE) return null;
                int maxLines = (int) maxLinesValue;
                boolean selectable = in.uint8() != 0;
                int richTextLength = in.int32();
                if (richTextLength < 0 || richTextLength > 1048576 || richTextLength > in.remaining()) return null;
                byte[] richText = in.bytes(richTextLength);
                String placeholder = in.requiredString();
                if (in.remaining() < 7) return null;
                int inputKind = in.uint8();
                boolean secure = in.uint8() != 0, multiline = in.uint8() != 0;
                int maxLength = in.int32();
                if (in.remaining() < 4) return null;
                int styleLength = in.int32();
                if (styleLength < 0 || styleLength > 1048576 || styleLength > in.remaining()) return null;
                byte[] typedStyle = in.bytes(styleLength);
                measured.add(measureNativeControl(id, kind, text, imageSource, typedStyle, minWidth, maxWidth, minHeight, maxHeight, textWrap, textOverflow, maxLines, selectable, richText, placeholder, inputKind, secure, multiline, maxLength));
            }
            if (in.hasRemaining()) return null;
            int capacity = 6;
            for (NativeMeasurement item : measured) capacity += 20 + item.error.getBytes(StandardCharsets.UTF_8).length;
            ByteBuffer out = ByteBuffer.allocate(capacity).order(ByteOrder.LITTLE_ENDIAN);
            out.putShort((short) 2).putInt(measured.size());
            for (NativeMeasurement item : measured) {
                byte[] error = item.error.getBytes(StandardCharsets.UTF_8);
                out.putLong(item.id).putFloat(item.width).putFloat(item.height).putInt(error.length).put(error);
            }
            return out.array();
        } catch (Throwable error) {
            Log.e("GoNative", "Native measurement batch failed", error);
            return null;
        }
    }

    public NativeMeasurement measureNativeControl(long id, int kind, String text, String imageSource, byte[] typedStyle,
                                                  float minWidth, float maxWidth, float minHeight, float maxHeight, int textWrap,
                                                  int textOverflow, int maxLines, boolean selectable, byte[] richText, String placeholder,
                                                  int inputKind, boolean secure, boolean multiline, int maxLength) {
        try {
            View view = factory.makeView(kind, false);
            if (view instanceof TextView) {
                TextView label = (TextView) view;
                CharSequence rendered = factory.decodeRichText(richText, 0, false);
                label.setText(rendered == null ? text : rendered);
                label.setIncludeFontPadding(false);
                label.setSingleLine(textWrap == 2 || (view instanceof EditText && !multiline));
                label.setMaxLines(maxLines > 0 ? maxLines : Integer.MAX_VALUE);
                label.setEllipsize(textOverflow == 1 ? TextUtils.TruncateAt.START : textOverflow == 2 ? TextUtils.TruncateAt.MIDDLE : textOverflow == 3 ? TextUtils.TruncateAt.END : null);
                if (view instanceof Button) {
                    Button button = (Button) view;
                    button.setAllCaps(false);
                    button.setMinWidth(0);
                    button.setMinimumWidth(0);
                    button.setMinHeight(factory.dp(44));
                    button.setMinimumHeight(factory.dp(44));
                    button.setPadding(factory.dp(16), 0, factory.dp(16), 0);
                }
                if (view instanceof EditText) {
                    EditText field = (EditText) view;
                    field.setSingleLine(!multiline);
                    field.setHint(placeholder);
                    field.setInputType(factory.resolveInputType(inputKind, 0, 0, secure, multiline));
                    field.setTransformationMethod(secure ? PasswordTransformationMethod.getInstance() : null);
                    field.setFilters(maxLength > 0 ? new InputFilter[]{new InputFilter.LengthFilter(maxLength)} : new InputFilter[0]);
                    field.setMinWidth(factory.dp(240));
                    field.setMinimumWidth(factory.dp(240));
                    field.setMinHeight(factory.dp(44));
                    field.setMinimumHeight(factory.dp(44));
                    field.setPadding(factory.dp(12), factory.dp(8), factory.dp(12), factory.dp(8));
                }
            }
            if (view instanceof ImageView && imageSource != null && !imageSource.isEmpty()) {
                int resource = context.getResources().getIdentifier(imageSource, "drawable", context.getPackageName());
                if (resource == 0) resource = context.getResources().getIdentifier(imageSource, "mipmap", context.getPackageName());
                if (resource != 0) ((ImageView) view).setImageResource(resource);
            }
            factory.applyTypedStyle(view, typedStyle);
            int widthSpec = nativeMeasureSpec(maxWidth);
            int heightSpec = nativeMeasureSpec(maxHeight);
            view.measure(widthSpec, heightSpec);
            float density = context.getResources().getDisplayMetrics().density;
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
        return View.MeasureSpec.makeMeasureSpec(factory.dp(maximum), View.MeasureSpec.AT_MOST);
    }

    private static boolean isFinite(float value) {
        return !Float.isNaN(value) && !Float.isInfinite(value);
    }

    private static float finiteNonNegative(float value) {
        return isFinite(value) ? Math.max(0, value) : 0;
    }
}
