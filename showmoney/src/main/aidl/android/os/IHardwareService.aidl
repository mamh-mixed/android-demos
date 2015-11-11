package android.os;

/** {@hide} */
interface IHardwareService
{
    // obsolete flashlight support
    boolean getFlashlightEnabled();
    void setFlashlightEnabled(boolean on);
}