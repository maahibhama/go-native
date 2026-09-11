// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "GoNativeKit",
    platforms: [.iOS(.v15)],
    products: [
        .library(name: "GoNativeKit", targets: ["GoNativeKit"]),
    ],
    targets: [
        .binaryTarget(
            name: "GoNativeKit",
            path: "build/native/GoNativeKit.xcframework"
        ),
    ]
)
