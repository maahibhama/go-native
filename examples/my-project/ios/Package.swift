// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "GoNativeAppDependencies",
    platforms: [.iOS(.v15)],
    products: [
        .library(name: "GoNativeAppDependencies", targets: ["GoNativeKit"]),
    ],
    targets: [
        .binaryTarget(name: "GoNativeKit", path: ".gonative/GoNativeKit.xcframework"),
    ]
)
