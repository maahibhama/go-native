package dev.gonative.runtime;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;

/** Bounded little-endian reader shared by Go Native's Android protocol decoders. */
public final class ProtocolReader {
    public static final int MAX_FIELD_BYTES = 1024 * 1024;

    private final ByteBuffer input;

    public ProtocolReader(byte[] payload) {
        if (payload == null) throw new IllegalArgumentException("payload is null");
        input = ByteBuffer.wrap(payload).order(ByteOrder.LITTLE_ENDIAN);
    }

    public int remaining() { return input.remaining(); }
    public boolean hasRemaining() { return input.hasRemaining(); }

    public int uint8() {
        require(1);
        return Byte.toUnsignedInt(input.get());
    }

    public int uint16() {
        require(2);
        return Short.toUnsignedInt(input.getShort());
    }

    public int int32() {
        require(4);
        return input.getInt();
    }

    public long uint32() { return Integer.toUnsignedLong(int32()); }

    public long int64() {
        require(8);
        return input.getLong();
    }

    public float float32() {
        require(4);
        return input.getFloat();
    }

    public byte[] bytes(int length) { return bytes(length, MAX_FIELD_BYTES); }

    public byte[] bytes(int length, int maximum) {
        if (length < 0 || length > maximum || length > remaining()) {
            throw new IllegalArgumentException("invalid byte field length");
        }
        byte[] result = new byte[length];
        input.get(result);
        return result;
    }

    public String requiredString() {
        return new String(bytes(int32()), StandardCharsets.UTF_8);
    }

    private void require(int count) {
        if (remaining() < count) throw new IllegalArgumentException("truncated protocol payload");
    }
}
