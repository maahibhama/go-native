package dev.gonative.runtime;

import android.animation.Animator;
import android.animation.AnimatorSet;
import android.animation.ObjectAnimator;
import android.animation.ValueAnimator;
import android.content.Context;
import android.graphics.Typeface;
import android.graphics.drawable.GradientDrawable;
import android.os.Build;
import android.text.Editable;
import android.text.InputFilter;
import android.text.InputType;
import android.text.SpannableStringBuilder;
import android.text.Spanned;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.text.method.LinkMovementMethod;
import android.text.method.PasswordTransformationMethod;
import android.text.style.ClickableSpan;
import android.text.style.ForegroundColorSpan;
import android.text.style.StyleSpan;
import android.text.style.UnderlineSpan;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.view.accessibility.AccessibilityEvent;
import android.view.accessibility.AccessibilityNodeInfo;
import android.view.animation.AccelerateDecelerateInterpolator;
import android.view.animation.AccelerateInterpolator;
import android.view.animation.DecelerateInterpolator;
import android.view.animation.LinearInterpolator;
import android.view.animation.OvershootInterpolator;
import android.view.inputmethod.EditorInfo;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.EditText;
import android.widget.HorizontalScrollView;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.ScrollView;
import android.widget.Switch;
import android.widget.TextView;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;

public final class ControlFactory {
    public static final int TEXT = 2;
    public static final int BUTTON = 3;
    public static final int ROW = 4;
    public static final int COLUMN = 5;
    public static final int SAFE_AREA = 6;
    public static final int TEXT_INPUT = 7;
    public static final int SWITCH = 8;
    public static final int PROGRESS_INDICATOR = 9;
    public static final int IMAGE = 10;
    public static final int SCROLL_VIEW = 11;

    private final Context context;
    private final ViewRegistry registry;
    private final EventDispatcher dispatcher;

    public ControlFactory(Context context, ViewRegistry registry, EventDispatcher dispatcher) {
        this.context = context;
        this.registry = registry;
        this.dispatcher = dispatcher;
    }

    public static final class NativeEditText extends EditText {
        long selectionHandler;
        final EventDispatcher dispatcher;

        NativeEditText(Context context, EventDispatcher dispatcher) {
            super(context);
            this.dispatcher = dispatcher;
        }

        @Override
        protected void onSelectionChanged(int start, int end) {
            super.onSelectionChanged(start, end);
            if (selectionHandler != 0 && dispatcher != null) {
                dispatcher.dispatchSelectionEvent(selectionHandler, start, end);
            }
        }
    }

    public View makeView(int kind, boolean horizontal) {
        if (kind == TEXT) {
            TextView tv = new TextView(context);
            tv.setTextColor(android.graphics.Color.BLACK);
            return tv;
        }
        if (kind == BUTTON) return new Button(context);
        if (kind == TEXT_INPUT) return new NativeEditText(context, dispatcher);
        if (kind == SWITCH) return new Switch(context);
        if (kind == PROGRESS_INDICATOR) {
            ProgressBar bar = new ProgressBar(context, null, android.R.attr.progressBarStyleHorizontal);
            bar.setMax(10000);
            return bar;
        }
        if (kind == IMAGE) return new ImageView(context);
        if (kind == SCROLL_VIEW) {
            if (horizontal) {
                HorizontalScrollView hsv = new HorizontalScrollView(context);
                hsv.setFillViewport(true);
                return hsv;
            } else {
                ScrollView sv = new ScrollView(context);
                sv.setFillViewport(true);
                return sv;
            }
        }
        LinearLayout layout = new LinearLayout(context);
        layout.setOrientation(kind == ROW ? LinearLayout.HORIZONTAL : LinearLayout.VERTICAL);
        if (kind == SAFE_AREA) layout.setFitsSystemWindows(true);
        return layout;
    }

    public void style(final View view, int kind, String text, float width, float height, float padding,
                      float gap, int alignment, float fontSize, boolean bold, long handler, long changeHandler,
                      long toggleHandler, boolean checked, float progress, String accessibility, String hint,
                      final int role, boolean focused, boolean scalesText, String imageSource, int imageMode,
                      int textWrap, int textOverflow, int maxLines, boolean selectable, byte[] richText,
                      long linkHandler, String placeholder, int inputMode, int inputKind, int returnKey,
                      int capitalization, int autoCorrect, boolean secure, final boolean multiline,
                      boolean readOnly, final int validationState, final String errorText, int selectionStart,
                      int selectionEnd, int maxLength, long submitHandler, long selectionHandler) {
        if (kind == TEXT && view instanceof TextView) {
            TextView textView = (TextView) view;
            CharSequence rendered = decodeRichText(richText, linkHandler, scalesText);
            textView.setText(rendered == null ? text : rendered);
            if (fontSize > 0) textView.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize);
            textView.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            textView.setIncludeFontPadding(false);
            textView.setTextColor(android.graphics.Color.parseColor("#111111"));
            textView.setSingleLine(textWrap == 2);
            textView.setHorizontallyScrolling(textWrap == 2);
            textView.setMaxLines(maxLines > 0 ? maxLines : Integer.MAX_VALUE);
            textView.setEllipsize(textOverflow == 1 ? TextUtils.TruncateAt.START : textOverflow == 2 ? TextUtils.TruncateAt.MIDDLE : textOverflow == 3 ? TextUtils.TruncateAt.END : null);
            textView.setTextIsSelectable(selectable);
            textView.setLinksClickable(linkHandler != 0);
            textView.setMovementMethod(linkHandler != 0 ? LinkMovementMethod.getInstance() : null);
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
            GradientDrawable btnBg = new GradientDrawable();
            btnBg.setColor(android.graphics.Color.parseColor("#007AFF"));
            btnBg.setCornerRadius(dp(8));
            btn.setBackground(btnBg);
            btn.setTextColor(android.graphics.Color.WHITE);
            btn.setPadding(dp(16), 0, dp(16), 0);
            final long eventHandler = handler;
            if (eventHandler != 0) {
                view.setOnClickListener(new View.OnClickListener() {
                    @Override public void onClick(View clicked) { if (dispatcher != null) dispatcher.dispatchEvent(eventHandler); }
                });
            } else {
                view.setOnClickListener(null);
            }
        }
        if (view instanceof EditText) {
            EditText field = (EditText) view;
            if (field instanceof NativeEditText) ((NativeEditText) field).selectionHandler = 0;
            field.setSingleLine(!multiline);
            field.setGravity(multiline ? Gravity.TOP | Gravity.START : Gravity.CENTER_VERTICAL);
            field.setIncludeFontPadding(false);
            field.setMinHeight(dp(44));
            field.setHint(placeholder);
            field.setTextColor(android.graphics.Color.BLACK);
            field.setHintTextColor(android.graphics.Color.parseColor("#8E8E93"));
            GradientDrawable fieldBg = new GradientDrawable();
            fieldBg.setColor(android.graphics.Color.parseColor("#FAFAFC"));
            fieldBg.setCornerRadius(dp(8));
            fieldBg.setStroke(dp(1), android.graphics.Color.parseColor("#D1D1D6"));
            field.setBackground(fieldBg);
            int padX = dp(12), padY = dp(8);
            field.setPadding(padX, padY, padX, padY);
            if (registry != null) {
                TextWatcher existing = registry.getTextWatcher(field);
                if (existing != null) field.removeTextChangedListener(existing);
            }
            boolean initialize = registry == null || !registry.isInputInitialized(field);
            if ((inputMode == 0 || initialize) && text != null && !field.getText().toString().equals(text)) field.setText(text);
            if (registry != null) registry.markInputInitialized(field);
            field.setTextSize(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, fontSize > 0 ? fontSize : 16);
            field.setTypeface(Typeface.DEFAULT, bold ? Typeface.BOLD : Typeface.NORMAL);
            field.setInputType(resolveInputType(inputKind, capitalization, autoCorrect, secure, multiline));
            field.setTransformationMethod(secure ? PasswordTransformationMethod.getInstance() : null);
            field.setImeOptions(resolveImeAction(returnKey) | (multiline && returnKey == 0 ? EditorInfo.IME_FLAG_NO_ENTER_ACTION : 0));
            field.setMaxLines(maxLines > 0 ? maxLines : (multiline ? Integer.MAX_VALUE : 1));
            field.setFilters(maxLength > 0 ? new InputFilter[]{new InputFilter.LengthFilter(maxLength)} : new InputFilter[0]);
            field.setFocusable(!readOnly);
            field.setFocusableInTouchMode(!readOnly);
            field.setCursorVisible(!readOnly);
            field.setLongClickable(!readOnly);
            if (validationState == 2) field.setError(errorText.isEmpty() ? "Invalid value" : errorText); else field.setError(null);
            final long actionHandler = submitHandler;
            field.setOnEditorActionListener(new TextView.OnEditorActionListener() {
                @Override public boolean onEditorAction(TextView ignored, int actionId, android.view.KeyEvent event) {
                    boolean enter = event != null && event.getKeyCode() == android.view.KeyEvent.KEYCODE_ENTER && event.getAction() == android.view.KeyEvent.ACTION_UP;
                    if (actionHandler != 0 && (actionId != EditorInfo.IME_ACTION_NONE || (!multiline && enter))) {
                        if (dispatcher != null) dispatcher.dispatchEvent(actionHandler);
                        return true;
                    }
                    return false;
                }
            });
            if (selectionStart >= 0 && selectionEnd >= selectionStart) {
                int safeStart = Math.min(selectionStart, field.length()), safeEnd = Math.min(selectionEnd, field.length());
                if (field.getSelectionStart() != safeStart || field.getSelectionEnd() != safeEnd) field.setSelection(safeStart, safeEnd);
            } else if (initialize) field.setSelection(field.length());
            if (field instanceof NativeEditText) ((NativeEditText) field).selectionHandler = selectionHandler;
            final long eventHandler = changeHandler;
            TextWatcher watcher = new TextWatcher() {
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                public void onTextChanged(CharSequence s, int start, int before, int count) {
                    if (eventHandler != 0 && dispatcher != null) dispatcher.dispatchValueEvent(eventHandler, s.toString());
                }
                public void afterTextChanged(Editable s) {}
            };
            field.addTextChangedListener(watcher);
            if (registry != null) registry.putTextWatcher(field, watcher);
        }
        if (view instanceof Switch) {
            Switch toggle = (Switch) view;
            toggle.setOnCheckedChangeListener(null);
            toggle.setChecked(checked);
            final long eventHandler = toggleHandler;
            if (eventHandler != 0) {
                toggle.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
                    @Override public void onCheckedChanged(CompoundButton button, boolean value) {
                        if (dispatcher != null) dispatcher.dispatchBoolEvent(eventHandler, value);
                    }
                });
            }
        }
        if (view instanceof ProgressBar) {
            ((ProgressBar) view).setProgress(Math.round(progress * 10000));
        }
        if (view instanceof ImageView) {
            ImageView image = (ImageView) view;
            if (imageSource != null && !imageSource.isEmpty()) {
                int resource = context.getResources().getIdentifier(imageSource, "drawable", context.getPackageName());
                if (resource == 0) resource = context.getResources().getIdentifier(imageSource, "mipmap", context.getPackageName());
                if (resource != 0) {
                    image.setImageResource(resource);
                } else if ("app_logo".equals(imageSource)) {
                    GradientDrawable gd = new GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setCornerRadius(dp(16));
                    image.setBackground(gd);
                    image.setImageResource(android.R.drawable.ic_lock_lock);
                    image.setColorFilter(android.graphics.Color.WHITE);
                    int pad = dp(12);
                    image.setPadding(pad, pad, pad, pad);
                } else if ("avatar".equals(imageSource)) {
                    GradientDrawable gd = new GradientDrawable();
                    gd.setColor(android.graphics.Color.parseColor("#007AFF"));
                    gd.setShape(GradientDrawable.OVAL);
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
        final String effectiveHint = hint;
        view.setAccessibilityDelegate(new View.AccessibilityDelegate() {
            @Override public void onInitializeAccessibilityNodeInfo(View host, AccessibilityNodeInfo info) {
                super.onInitializeAccessibilityNodeInfo(host, info);
                if (role == 1 || role == 3) info.setClassName(TextView.class.getName());
                else if (role == 2) info.setClassName(Button.class.getName());
                else if (role == 4) info.setClassName("android.widget.ImageView");
                if (Build.VERSION.SDK_INT >= 26 && !effectiveHint.isEmpty()) info.setHintText(effectiveHint);
                if (Build.VERSION.SDK_INT >= 28 && role == 3) info.setHeading(true);
                if (view instanceof EditText && validationState == 2) {
                    info.setContentInvalid(true);
                    info.setError(errorText.isEmpty() ? "Invalid value" : errorText);
                }
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

    public void applyComputedFrame(long nodeID, final View view, boolean hasFrame, float x, float y, float width, float height, boolean isRoot) {
        if (!hasFrame || view == null || isRoot) return;
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

    public void applyTypedStyle(View view, byte[] payload) {
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
            String family = new String(payload, textBase + 185, fontLength, StandardCharsets.UTF_8);
            Typeface face = family.isEmpty() ? Typeface.DEFAULT : Typeface.create(family, Typeface.NORMAL);
            text.setTypeface(face, fontWeight >= 600 ? Typeface.BOLD : Typeface.NORMAL);
            if (fontSize > 0) text.setTextSize(TypedValue.COMPLEX_UNIT_DIP, fontSize);
            if (lineHeight > 0 && Build.VERSION.SDK_INT >= 28) text.setLineHeight(dp(lineHeight));
            if (letterSpacing != 0 && fontSize > 0) text.setLetterSpacing(letterSpacing / fontSize);
        }
        if (android.graphics.Color.alpha(background) > 0 || borderWidth > 0) {
            GradientDrawable drawable = new GradientDrawable();
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

    public void applyInteractions(long nodeID, View view, byte[] payload) {
        if (view == null) return;
        if (registry != null) {
            GestureBinding old = registry.getGestureBinding(nodeID);
            if (old != null) old.dispose();
            registry.removeGestureBinding(nodeID);
        }
        if (payload == null || payload.length < 4) return;
        ProtocolReader in = new ProtocolReader(payload);
        ArrayList<GestureBinding.GestureSpec> gestures = new ArrayList<>();
        int gestureCount = in.int32();
        for (int i = 0; i < gestureCount && in.remaining() >= 26; i++) {
            GestureBinding.GestureSpec spec = new GestureBinding.GestureSpec();
            spec.kind = in.uint8();
            spec.direction = in.uint8();
            spec.minimumPressNanos = in.int64();
            spec.minimumTravel = in.float32();
            spec.handler = in.int64();
            gestures.add(spec);
        }
        if (!gestures.isEmpty()) {
            GestureBinding binding = new GestureBinding(view, gestures, dispatcher);
            if (registry != null) registry.putGestureBinding(nodeID, binding);
            view.setOnTouchListener(binding);
        } else view.setOnTouchListener(null);

        if (in.remaining() < 4) return;
        int animationCount = in.int32();
        ArrayList<Animator> animations = new ArrayList<>();
        for (int i = 0; i < animationCount && in.remaining() >= 42; i++) {
            int property = in.uint8();
            long durationNanos = in.int64();
            long delayNanos = in.int64();
            int curve = in.uint8();
            float damping = in.float32();
            float velocity = in.float32();
            boolean reduceMotionOK = in.uint8() != 0;
            float from = in.float32(), to = in.float32(), fromX = in.float32(), fromY = in.float32(), toX = in.float32(), toY = in.float32();
            Animator animator = makeAnimator(view, property, from, to, fromX, fromY, toX, toY);
            if (animator == null) continue;
            animator.setDuration(Math.max(0, durationNanos / 1000000L));
            animator.setStartDelay(Math.max(0, delayNanos / 1000000L));
            if (curve == 1) animator.setInterpolator(new AccelerateInterpolator());
            else if (curve == 2) animator.setInterpolator(new DecelerateInterpolator());
            else if (curve == 3) animator.setInterpolator(new LinearInterpolator());
            else if (curve == 4) animator.setInterpolator(new OvershootInterpolator(Math.max(.1f, (1f - damping) * 2f + Math.abs(velocity) * .1f)));
            else animator.setInterpolator(new AccelerateDecelerateInterpolator());
            if (!reduceMotionOK && Build.VERSION.SDK_INT >= 26 && !ValueAnimator.areAnimatorsEnabled()) animator.setDuration(0);
            animations.add(animator);
        }
        if (!animations.isEmpty()) { AnimatorSet set = new AnimatorSet(); set.playSequentially(animations); set.start(); }
    }

    public CharSequence decodeRichText(byte[] payload, final long linkHandler, boolean scalesText) {
        if (payload == null || payload.length == 0) return null;
        try {
            ProtocolReader in = new ProtocolReader(payload);
            if (in.remaining() < 6 || in.uint16() != 1) return null;
            long count = in.uint32();
            if (count > 100000) return null;
            SpannableStringBuilder output = new SpannableStringBuilder();
            for (long i = 0; i < count; i++) {
                final String spanText = in.requiredString(), link = in.requiredString();
                if (in.remaining() < 12) return null;
                float size = in.float32(); int weight = in.uint16(); int flags = in.uint8(); boolean hasColor = in.uint8() != 0;
                int red = in.uint8(), green = in.uint8(), blue = in.uint8(), alpha = in.uint8();
                int color = android.graphics.Color.argb(alpha, red, green, blue);
                int start = output.length(); output.append(spanText); int end = output.length();
                if (size > 0) {
                    int pixels = Math.round(TypedValue.applyDimension(scalesText ? TypedValue.COMPLEX_UNIT_SP : TypedValue.COMPLEX_UNIT_DIP, size, context.getResources().getDisplayMetrics()));
                    output.setSpan(new android.text.style.AbsoluteSizeSpan(pixels), start, end, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
                }
                if (weight >= 600 || (flags & 2) != 0) output.setSpan(new StyleSpan(weight >= 600 && (flags & 2) != 0 ? Typeface.BOLD_ITALIC : weight >= 600 ? Typeface.BOLD : Typeface.ITALIC), start, end, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
                if ((flags & 1) != 0) output.setSpan(new UnderlineSpan(), start, end, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
                if (hasColor) output.setSpan(new ForegroundColorSpan(color), start, end, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
                if (!link.isEmpty() && linkHandler != 0) output.setSpan(new ClickableSpan() {
                    @Override public void onClick(View widget) { if (dispatcher != null) dispatcher.dispatchValueEvent(linkHandler, link); }
                }, start, end, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
            }
            return in.hasRemaining() ? null : output;
        } catch (Throwable ignored) { return null; }
    }

    public int dp(float value) {
        return Math.round(value * context.getResources().getDisplayMetrics().density);
    }

    public int resolveInputType(int kind, int capitalization, int autoCorrect, boolean secure, boolean multiline) {
        int type;
        if (kind == 1) type = InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_EMAIL_ADDRESS;
        else if (kind == 2) type = InputType.TYPE_CLASS_PHONE;
        else if (kind == 3) type = InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_URI;
        else if (kind == 4) type = InputType.TYPE_CLASS_NUMBER | InputType.TYPE_NUMBER_FLAG_SIGNED;
        else if (kind == 5) type = InputType.TYPE_CLASS_NUMBER | InputType.TYPE_NUMBER_FLAG_DECIMAL | InputType.TYPE_NUMBER_FLAG_SIGNED;
        else type = InputType.TYPE_CLASS_TEXT | (kind == 6 ? InputType.TYPE_TEXT_VARIATION_FILTER : InputType.TYPE_TEXT_VARIATION_NORMAL);
        if ((type & InputType.TYPE_CLASS_TEXT) != 0) {
            if (capitalization == 1) type |= InputType.TYPE_TEXT_FLAG_CAP_SENTENCES;
            else if (capitalization == 2) type |= InputType.TYPE_TEXT_FLAG_CAP_WORDS;
            else if (capitalization == 3) type |= InputType.TYPE_TEXT_FLAG_CAP_CHARACTERS;
            if (autoCorrect == 1) type |= InputType.TYPE_TEXT_FLAG_AUTO_CORRECT;
            else if (autoCorrect == 2) type |= InputType.TYPE_TEXT_FLAG_NO_SUGGESTIONS;
            if (multiline) type |= InputType.TYPE_TEXT_FLAG_MULTI_LINE;
            if (secure) type = (type & ~InputType.TYPE_MASK_VARIATION) | InputType.TYPE_TEXT_VARIATION_PASSWORD;
        }
        return type;
    }

    public int resolveImeAction(int key) {
        if (key == 1) return EditorInfo.IME_ACTION_DONE;
        if (key == 2) return EditorInfo.IME_ACTION_GO;
        if (key == 3) return EditorInfo.IME_ACTION_NEXT;
        if (key == 4) return EditorInfo.IME_ACTION_SEARCH;
        if (key == 5) return EditorInfo.IME_ACTION_SEND;
        return EditorInfo.IME_ACTION_NONE;
    }

    private Animator makeAnimator(final View view, int property, float from, float to, float fromX, float fromY, float toX, float toY) {
        if (property == 1) return ObjectAnimator.ofFloat(view, View.ALPHA, from, to);
        if (property == 2) {
            AnimatorSet set = new AnimatorSet();
            set.playTogether(ObjectAnimator.ofFloat(view, View.SCALE_X, from, to), ObjectAnimator.ofFloat(view, View.SCALE_Y, from, to));
            return set;
        }
        if (property == 3) {
            AnimatorSet set = new AnimatorSet();
            set.playTogether(ObjectAnimator.ofFloat(view, View.TRANSLATION_X, dp(fromX), dp(toX)), ObjectAnimator.ofFloat(view, View.TRANSLATION_Y, dp(fromY), dp(toY)));
            return set;
        }
        if (property == 4) {
            ValueAnimator animator = ValueAnimator.ofFloat(0, 1);
            animator.addUpdateListener(new ValueAnimator.AnimatorUpdateListener() {
                @Override public void onAnimationUpdate(ValueAnimator animation) { view.requestLayout(); }
            });
            return animator;
        }
        return null;
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

    private static boolean isFinite(float value) {
        return !Float.isNaN(value) && !Float.isInfinite(value);
    }
}
