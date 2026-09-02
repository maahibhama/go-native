//go:build android

#include <jni.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>

extern void GoNativeAndroidStart(void);
extern void GoNativeAndroidDispatchEvent(uint64_t handler);
extern void GoNativeAndroidStop(void);
extern void GoNativeAndroidReportBatchApplied(uint64_t sequence, uint64_t nativeNanos);

static JavaVM *gn_vm;
static jobject gn_renderer;
static jmethodID gn_apply;
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
    (*env)->DeleteLocalRef(env, cls);
    pthread_mutex_unlock(&gn_renderer_mu);
    GoNativeAndroidStart();
}

JNIEXPORT void JNICALL
Java_dev_gonative_counter_MainActivity_nativeDispatchEvent(JNIEnv *env, jobject renderer, jlong handler) {
    (void)env;
    (void)renderer;
    GoNativeAndroidDispatchEvent((uint64_t)handler);
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
    pthread_mutex_unlock(&gn_renderer_mu);
}

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
