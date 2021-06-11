
## Running Simulator ##
Starling simulator has a single executable. This contains the backend server as well as the UX website bundled together.

Platform      | Executable
--------------|----------------------------------
Windows       | `bin/starling_windows_amd64.exe`
macOS         | `bin/starling_darwin_amd64`
Linux         | `bin/starling_linux_amd64`
Raspberry Pi  | `bin/starling_linux_arm64`

### Running Simulation Server ###
To start the starling simulation server, run the above executable.

<img src="assets/start.png" alt="Starting Starling" height=150 />

### Create Central application ###
Starling simulates devices that connect to an IoT Central application.
So, create an IoT Central application and create your device templates.
You can use one of the device model samples: [brewer.json](./brewer.json) or [drone.json](./drone.json) .
After configuring the views, publish these device templates. Create an API Token with administrator role and copy it.

[Back to contents](../README.md)| Previous: [Building binaries](build.md) | Next: [Configuring and running simulations](configure.md)
---------------------------------|-------------------------------------------------------|------------------------------------
