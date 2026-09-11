package dev.gonative.runtime;

public interface EventDispatcher {
    void dispatchEvent(long handler);
    void dispatchValueEvent(long handler, String value);
    void dispatchBoolEvent(long handler, boolean value);
    void dispatchGestureEvent(long handler, float translationX, float translationY, float velocityX, float velocityY);
    void dispatchSelectionEvent(long handler, int start, int end);
    void dispatchFocus(long nodeID, boolean focused);
}
