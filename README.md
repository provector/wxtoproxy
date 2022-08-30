# wxtoproxy

This is a simple proxy "wrapper" for the ageing [WxToImage](https://wxtoimgrestored.xyz/). It starts a proxy server on 127.0.0.1 port 8080 and properly handles any application redirects. When the weather.txt format breaks in future its possible to implement a data converter here as well. 

## Setup

Before first run (or during first run) make sure to set up proxy setting in WxToImg Options->Internet Options, check proxy set host to ```127.0.0.1``` and port to ```8080```

![image](https://user-images.githubusercontent.com/12935423/177414352-1ec1ecac-ff5b-410a-b98b-58b41eff2769.png)

Copy the ```wxtoproxy.exe``` file to your WxToImg directory. Double-click and enjoy. Updating Kepplers should work fine, you can see any input and application calls in the terminal window. Program should automatically close when you close the app. Closing terminal window will most likely stop updating to work :)

If you are running for some reason 32 bit system (windows 7) there is a 32 bit compiled version as well: ```wxtoproxy_win7_32.exe```

Credit for proxy code: [yowu](https://gist.github.com/yowu/f7dc34bd4736a65ff28d)


## Linux

For Linux users a simple "wget" command will do the trick (sudo might be required)

  wget -O /usr/local/lib/wx/tle/weather.txt http://celestrak.org/NORAD/elements/weather.txt
  
This command will replace the expired weather.txt file. 
