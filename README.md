This program is intended to generate and configure SMF service manifests on Oracle Solaris or Illumos based systems.
While this program can be run on any other operating system that supports the project's dependecies, manifest valiadation support is
only available on Oracle Solaris or Illumos based systems with the 'svccfg' binary.

This program comes with NO SUPPORT and is not intended for deployment in production environments. Use at your own risk.

**Building**: Run the following commands to clone the repository, enter the directory, and invoke the build script.
<br />`git clone https://github.com/madelinehebert/MDHsvcbundle-go`
<br />`cd MDHsvcbundle-go`
<br />`bash ./build.sh`

If you are attempting to build this on Windows, just run the contents of the build script in your terminal.

**Installing pkg release on Illumos**: Run the following commands to download the package release and add it to your system.
<br /> `wget $LINK_TO_PKG_RELEASE -o $OUTPUT_LOCATION`
<br /> `pkgadd -d $PATH_TO_OUTPUT_LOCATION`

**Uninstalling pkg release on Illumos**: Run the following command to uninstall this package from your system. (This is assuming you have not moved any installed files.)
<br /> `pkgrm svcbundle`
