This is idiomatic simple example of using generic-ms-boot with nutz/up,
suitable as a copy-paste-worth boilerplate.

Running this program:

    $ task x1

    task: [build] go build -o x1.xu
    task: [x1] ./x1.xu ./config-01.yaml . SHORT-192 96bd92c9c7d2 
    {"level":"INFO","PRODUCT":"Idiomatic Simple Example by Someone","COPYRIGHT":"(c) Someone, 2025","tags":["BOOT-STARTED"],"etc":[1],"tsz":"11.481µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["CONFIGDB-INIT","ziPdTJw"],"etc":[1,1],"msg":"reading the config file(s)...","tsz":"93.066µs","tsa":"20250709-191112"}
    {"level":"INFO","depth":-1,"tags":["CONFIGDB-INIT","-1","EZWLkX"],"etc":[1,2,1],"msg":"reading the config file './config-01.yaml'...","tsz":"100.61µs","tsa":"20250709-191112"}
    {"level":"INFO","depth":-1,"tags":["CONFIGDB-INIT","-1","KPPEY7ZW"],"etc":[1,2,2],"msg":"config file './config-01.yaml' id='' err='{: FORWARDTRACE {: ETC [1 3], TAGS [CONFIGDB-INIT]}, BACKTRACE { [e6269413b51b \"nonexistent map key 'id'\"]}}'","tsz":"242.339µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["CONFIGDB-INIT","K2aUDgz"],"etc":[1,9],"msg":"reading the config file(s) OK","tsz":"256.385µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","BEGIN"],"etc":[4,1],"tsz":"287.745µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","core","START"],"etc":[4,2,1,1],"tsz":"292.674µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","core","0198437cd198"],"etc":[4,2,1,2],"msg":"core--b154c4ee-9049-4ec7-aeb7-484ffaa4b7a4 overridden core--084433f7-4260-4ef8-9cee-fd36f5f1de56 12 true","tsz":"297.383µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys1","START"],"etc":[4,2,2,1],"tsz":"310.538µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys1","3dde2aad139c"],"etc":[4,2,2,2],"msg":"sys1--dea14a4f-a99a-40e7-bac9-3b15d1283fa6 overridden sys1--acdf81f5-1942-4358-a677-60d0b8bd688f 12 true","tsz":"313.975µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys2","START"],"etc":[4,2,3,1],"tsz":"320.948µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys2","9b971a886a1a"],"etc":[4,2,3,2],"msg":"sys2--ec1b0a6f-81f2-4901-bdea-39913b2dd6ce overridden sys2--07fd2176-247f-4b6d-9039-f8a121cd363d 12 true","tsz":"325.046µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","START"],"etc":[4,2,4,1],"tsz":"330.927µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","0198437cd198"],"etc":[4,2,4,2],"msg":"9.87654321e+08 hello 12 true","tsz":"337.96µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","52"],"etc":[4,2,4,3],"msg":"DiagnosticFlag 96bd92c9c7d2 activated","tsz":"342.238µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","UP"],"etc":[4,3],"tsz":"346.326µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","BAW-227"],"etc":[4,4],"msg":"block and wait","tsz":"351.015µs","tsa":"20250709-191112"}
    {"level":"INFO","tags":["mwt","BAW-235"],"etc":[4,5],"msg":"got system interrupt signal","tsz":"2.535457553s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","GO-DOWN"],"etc":[4,6],"tsz":"2.535469817s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","RWB-33"],"etc":[4,7,1],"msg":"rec-windup-begin","tsz":"2.535478994s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Service","STOP"],"etc":[4,7,2,1],"tsz":"2.535483573s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Sys2","STOP"],"etc":[4,7,3,1],"tsz":"2.535487891s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Sys1","STOP"],"etc":[4,7,4,1],"tsz":"2.535493562s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","core","STOP"],"etc":[4,7,5,1],"tsz":"2.535497479s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","RWE-34"],"etc":[4,7,6],"msg":"rec-windup-end","tsz":"2.535501627s","tsa":"20250709-191115"}
    {"level":"INFO","tags":["mwt","END"],"etc":[4,8],"tsz":"2.535506246s","tsa":"20250709-191115"}

