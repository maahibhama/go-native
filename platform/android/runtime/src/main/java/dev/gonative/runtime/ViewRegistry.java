package dev.gonative.runtime;

import android.text.TextWatcher;
import android.util.LongSparseArray;
import android.view.View;
import android.view.ViewGroup;
import android.widget.EditText;

import java.util.WeakHashMap;

public final class ViewRegistry {
    private final LongSparseArray<View> views = new LongSparseArray<>();
    private final LongSparseArray<GestureBinding> gestureBindings = new LongSparseArray<>();
    private final WeakHashMap<EditText, TextWatcher> textWatchers = new WeakHashMap<>();
    private final WeakHashMap<EditText, Boolean> initializedInputs = new WeakHashMap<>();
    private long rootNodeID;

    public View getView(long id) {
        return views.get(id);
    }

    public void putView(long id, View view) {
        views.put(id, view);
    }

    public int viewCount() {
        return views.size();
    }

    public View viewAt(int index) {
        return views.valueAt(index);
    }

    public long getRootNodeID() {
        return rootNodeID;
    }

    public void setRootNodeID(long id) {
        this.rootNodeID = id;
    }

    public GestureBinding getGestureBinding(long id) {
        return gestureBindings.get(id);
    }

    public void putGestureBinding(long id, GestureBinding binding) {
        gestureBindings.put(id, binding);
    }

    public GestureBinding removeGestureBinding(long id) {
        GestureBinding binding = gestureBindings.get(id);
        gestureBindings.remove(id);
        return binding;
    }

    public TextWatcher getTextWatcher(EditText field) {
        return textWatchers.get(field);
    }

    public void putTextWatcher(EditText field, TextWatcher watcher) {
        textWatchers.put(field, watcher);
    }

    public void removeTextWatcher(EditText field) {
        textWatchers.remove(field);
    }

    public boolean isInputInitialized(EditText field) {
        return initializedInputs.containsKey(field);
    }

    public void markInputInitialized(EditText field) {
        initializedInputs.put(field, Boolean.TRUE);
    }

    public static void detach(View view) {
        if (view != null && view.getParent() instanceof ViewGroup) {
            ((ViewGroup) view.getParent()).removeView(view);
        }
    }

    public void deleteNode(long nodeID) {
        View view = views.get(nodeID);
        detach(view);
        if (view instanceof EditText) {
            textWatchers.remove((EditText) view);
            initializedInputs.remove((EditText) view);
        }
        GestureBinding binding = gestureBindings.get(nodeID);
        if (binding != null) {
            binding.dispose();
        }
        gestureBindings.remove(nodeID);
        if (view != null) {
            view.animate().cancel();
        }
        views.remove(nodeID);
        if (nodeID == rootNodeID) {
            rootNodeID = 0;
        }
    }

    public void clear() {
        for (int i = 0; i < gestureBindings.size(); i++) {
            gestureBindings.valueAt(i).dispose();
        }
        gestureBindings.clear();
        textWatchers.clear();
        initializedInputs.clear();
        for (int i = 0; i < views.size(); i++) {
            views.valueAt(i).animate().cancel();
        }
        views.clear();
        rootNodeID = 0;
    }
}
