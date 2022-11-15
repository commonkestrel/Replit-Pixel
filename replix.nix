{ pkgs }: {
    deps = [
        pkgs.busybox
        pkgs.go_1_18
        pkgs.gopls

        #Needed for glfw
        pkgs.glib
        pkgs.glfw
        pkgs.xorg.libX11
        pkgs.xorg.libXext
        pkgs.xorg.libXinerama
        pkgs.xorg.libXcursor
        pkgs.xorg.libXrandr
        pkgs.xorg.libXi
        pkgs.xorg.libXxf86vm
        #Needed for pixel
        pkgs.pkgconfig
    ];
}
