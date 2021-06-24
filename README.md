# busybox64.portable
Just a barebones busybox environment for you to expand upon, 

Renaming the EXE renames the datafolder as well, so if you have previously ran the program as 
`busybox64.portable.exe`, and you rename it to `shell.exe`, you will need to rename your datafolder 
to `shell.data`

# Credits
This program utilizes a [Prebuilt version of Busybox by Frippery](https://frippery.org/busybox/index.html) to support basic shell environments. Using this you can setup your own portable *unix-like* shell however you like.

If you need it to just run a single shell script such as [Weeb2PSP](https://github.com/tuxlovesyou/weeb2psp), You can configure it to run a shell script to fetch files and then run the actual commands needed.

# Tools
[[PSPConverter](https://mega.nz/folder/PsIEGJzJ#0h9c-bIboeNWgVpk1u2cGQ)]: To use, extract to the `VIDEOS` folder on your psp MEMSTICK, and then put the videos you wish to convert for your PSP, into the `video_in` folder, this can be changed by modifying `config.toml` to change the `INPUT` and `OUTPUT` folder destinations
