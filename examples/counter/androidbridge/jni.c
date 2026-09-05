//go:build android

#include <jni.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>

extern void GoNativeAndroidStart(void);
extern void GoNativeAndroidDispatchEvent(uint64_t handler);
extern void GoNativeAndroidDispatchValueEvent(uint64_t handler, const char *value);
extern void GoNativeAndroidDispatchBoolEvent(uint64_t handler, uint8_t value);
extern void GoNativeAndroidDispatchGestureEvent(uint64_t handler, float translationX, float translationY, float velocityX, float velocityY);
extern void GoNativeAndroidStop(void);
extern void GoNativeAndroidReportBatchApplied(uint64_t sequence, uint64_t nativeNanos);
extern void GoNativeAndroidSetLifecycle(uint8_t state);
extern void GoNativeAndroidDispatchFocus(uint64_t nodeID, uint8_t focused);
extern void GoNativeAndroidUpdateViewport(float width, float height, float scale);

static JavaVM *gn_vm;
static jobject gn_renderer;
static jmethodID gn_apply;
static jmethodID gn_measure;
static pthread_mutex_t gn_renderer_mu = PTHREAD_MUTEX_INITIALIZER;

JNIEXPORT jint JNICALL JNI_OnLoad(JavaVM *vm, void *reserved) {
    (void)reserved;
    gn_vm = vm;
    return JNI_VERSION_1_6;
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeStart(JNIEnv *env, jobject renderer) {
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer) {
        (*env)->DeleteGlobalRef(env, gn_renderer);
    }
    gn_renderer = (*env)->NewGlobalRef(env, renderer);
    jclass cls = (*env)->GetObjectClass(env, renderer);
    gn_apply = (*env)->GetMethodID(env, cls, "applyMutationBatch", "([B)V");
    gn_measure = (*env)->GetMethodID(env, cls, "measureNativeBatch", "([B)[B");
    (*env)->DeleteLocalRef(env, cls);
    pthread_mutex_unlock(&gn_renderer_mu);
    GoNativeAndroidStart();
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeSetLifecycle(JNIEnv *env, jobject renderer, jint state) {
    (void)env;
    (void)renderer;
    if (state >= 0 && state <= 6) GoNativeAndroidSetLifecycle((uint8_t)state);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchFocus(JNIEnv *env, jobject renderer, jlong nodeID, jboolean focused) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchFocus((uint64_t)nodeID, focused ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeUpdateViewport(JNIEnv *env, jobject renderer, jfloat width, jfloat height, jfloat scale) {
    (void)env;
    (void)renderer;
    GoNativeAndroidUpdateViewport((float)width, (float)height, (float)scale);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchEvent(JNIEnv *env, jobject renderer, jlong handler) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchEvent((uint64_t)handler);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchValueEvent(JNIEnv *env, jobject renderer, jlong handler, jstring value) {
    (void)renderer;
    const char *utf8 = value ? (*env)->GetStringUTFChars(env, value, NULL) : "";
    if (utf8) { GoNativeAndroidDispatchValueEvent((uint64_t)handler, utf8); }
    if (value && utf8) { (*env)->ReleaseStringUTFChars(env, value, utf8); }
}
JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchBoolEvent(JNIEnv *env, jobject renderer, jlong handler, jboolean value) {
    (void)env; (void)renderer; GoNativeAndroidDispatchBoolEvent((uint64_t)handler, value ? 1 : 0);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchGestureEvent(JNIEnv *env, jobject renderer, jlong handler, jfloat translationX, jfloat translationY, jfloat velocityX, jfloat velocityY) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchGestureEvent((uint64_t)handler, (float)translationX, (float)translationY, (float)velocityX, (float)velocityY);
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeStop(JNIEnv *env, jobject renderer) {
    (void)renderer;
    GoNativeAndroidStop();
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer) {
        (*env)->DeleteGlobalRef(env, gn_renderer);
        gn_renderer = NULL;
    }
    gn_apply = NULL;
    gn_measure = NULL;
    pthread_mutex_unlock(&gn_renderer_mu);
}

int32_t GNAndroidMeasureBatch(const uint8_t *bytes, int32_t length, uint8_t **result) {
    if (result) *result = NULL;
    if (!gn_vm || !result || !bytes || length <= 0) return -1;
    JNIEnv *env = NULL;
    int attached = 0;
    jint status = (*gn_vm)->GetEnv(gn_vm, (void **)&env, JNI_VERSION_1_6);
    if (status == JNI_EDETACHED) {
        if ((*gn_vm)->AttachCurrentThread(gn_vm, &env, NULL) != JNI_OK) return -2;
        attached = 1;
    } else if (status != JNI_OK) return -2;
    int32_t result_length = -3;
    pthread_mutex_lock(&gn_renderer_mu);
    if (gn_renderer && gn_measure) {
        jbyteArray request = (*env)->NewByteArray(env, length);
        if (request) {
            (*env)->SetByteArrayRegion(env, request, 0, length, (const jbyte *)bytes);
            jbyteArray response = (jbyteArray)(*env)->CallObjectMethod(env, gn_renderer, gn_measure, request);
            if (!(*env)->ExceptionCheck(env) && response) {
                jsize response_length = (*env)->GetArrayLength(env, response);
                if (response_length > 0 && response_length <= 16777216) {
                    uint8_t *copy = (uint8_t *)malloc((size_t)response_length);
                    if (copy) {
                        (*env)->GetByteArrayRegion(env, response, 0, response_length, (jbyte *)copy);
                        *result = copy;
                        result_length = (int32_t)response_length;
                    } else result_length = -4;
                }
                (*env)->DeleteLocalRef(env, response);
            }
            (*env)->DeleteLocalRef(env, request);
        }
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
            result_length = -5;
        }
    }
    pthread_mutex_unlock(&gn_renderer_mu);
    if (attached) (*gn_vm)->DetachCurrentThread(gn_vm);
    return result_length;
}

void GNAndroidFreeBuffer(uint8_t *bytes) { free(bytes); }

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeReportBatchApplied(JNIEnv *env, jobject renderer, jlong sequence, jlong nativeNanos) {
    (void)env;
    (void)renderer;
    GoNativeAndroidReportBatchApplied((uint64_t)sequence, (uint64_t)nativeNanos);
}

void GNAndroidApplyMutationBatch(const uint8_t *bytes, int32_t length) {
    if (!gn_vm || length <= 0) {
        return;
    }
    JNIEnv *env = NULL;
    int attached = 0;
    jint status = (*gn_vm)->GetEnv(gn_vm, (void **)&env, JNI_VERSION_1_6);
    if (status == JNI_EDETACHED) {
        if ((*gn_vm)->AttachCurrentThread(gn_vm, &env, NULL) != JNI_OK) {
            return;
        }
        attached = 1;
    } else if (status != JNI_OK) {
        return;
    }
    pthread_mutex_lock(&gn_renderer_mu);
    if (!gn_renderer || !gn_apply) {
        pthread_mutex_unlock(&gn_renderer_mu);
        if (attached) (*gn_vm)->DetachCurrentThread(gn_vm);
        return;
    }
    jbyteArray payload = (*env)->NewByteArray(env, length);
    if (payload) {
        (*env)->SetByteArrayRegion(env, payload, 0, length, (const jbyte *)bytes);
        (*env)->CallVoidMethod(env, gn_renderer, gn_apply, payload);
        (*env)->DeleteLocalRef(env, payload);
    }
    if ((*env)->ExceptionCheck(env)) {
        (*env)->ExceptionDescribe(env);
        (*env)->ExceptionClear(env);
    }
    pthread_mutex_unlock(&gn_renderer_mu);
    if (attached) {
        (*gn_vm)->DetachCurrentThread(gn_vm);
    }
}
