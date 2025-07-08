This is idiomatic simple example of using generic-ms-boot with nutz/up,
suitable as a copy-paste-worth boilerplate.

Running this program:

    $ task x1

    task: [build] go build -o x1.xu
    task: [x1] ./x1.xu ./config-01.yaml . SHORT-192 96bd92c9c7d2 
    {"level":"INFO","tags":["CONFIGDB-INIT","ziPdTJw"],"etc":[0,1],"msg":"reading the config file(s)...","tsz":"21.799µs","tsa":"20250709-001335"}
    {"level":"INFO","depth":-1,"tags":["CONFIGDB-INIT","-1","EZWLkX"],"etc":[0,2,1],"msg":"reading the config file './config-01.yaml'...","tsz":"179.021µs","tsa":"20250709-001335"}
    {"level":"INFO","depth":-1,"tags":["CONFIGDB-INIT","-1","KPPEY7ZW"],"etc":[0,2,2],"msg":"config file './config-01.yaml' id='' err='{: FORWARDTRACE {: ETC [0 3], TAGS [CONFIGDB-INIT]}, BACKTRACE { [e6269413b51b \"nonexistent map key 'id'\"]}}'","tsz":"417.308µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["CONFIGDB-INIT","K2aUDgz"],"etc":[0,9],"msg":"reading the config file(s) OK","tsz":"444.697µs","tsa":"20250709-001335"}
    {"level":"INFO","copyright":"(c) Dag Manaskent, 2025","tags":["mwt","BEGIN"],"etc":[3,1],"tsz":"502.741µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","core","START"],"etc":[3,2,1,1],"tsz":"513.48µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","core","0198437cd198"],"etc":[3,2,1,2],"msg":"core--b154c4ee-9049-4ec7-aeb7-484ffaa4b7a4 overridden core--084433f7-4260-4ef8-9cee-fd36f5f1de56 12 true","tsz":"522.166µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys1","START"],"etc":[3,2,2,1],"tsz":"548.273µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys1","3dde2aad139c"],"etc":[3,2,2,2],"msg":"sys1--dea14a4f-a99a-40e7-bac9-3b15d1283fa6 overridden sys1--acdf81f5-1942-4358-a677-60d0b8bd688f 12 true","tsz":"554.223µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys2","START"],"etc":[3,2,3,1],"tsz":"566.565µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Sys2","9b971a886a1a"],"etc":[3,2,3,2],"msg":"sys2--ec1b0a6f-81f2-4901-bdea-39913b2dd6ce overridden sys2--07fd2176-247f-4b6d-9039-f8a121cd363d 12 true","tsz":"571.494µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","START"],"etc":[3,2,4,1],"tsz":"583.686µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","0198437cd198"],"etc":[3,2,4,2],"msg":"9.87654321e+08 hello 12 true","tsz":"593.564µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","otdg-main","Service","52"],"etc":[3,2,4,3],"msg":"DiagnosticFlag 96bd92c9c7d2 activated","tsz":"599.144µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","UP"],"etc":[3,3],"tsz":"604.493µs","tsa":"20250709-001335"}
    {"level":"INFO","tags":["mwt","BAW-227"],"etc":[3,4],"msg":"block and wait","tsz":"609.192µs","tsa":"20250709-001335"}
    ^C
    {"level":"INFO","tags":["mwt","BAW-235"],"etc":[3,5],"msg":"got system interrupt signal","tsz":"15.20781868s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","GO-DOWN"],"etc":[3,6],"tsz":"15.207847792s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","RWB-33"],"etc":[3,7,1],"msg":"rec-windup-begin","tsz":"15.207882915s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Service","STOP"],"etc":[3,7,2,1],"tsz":"15.207897421s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Sys2","STOP"],"etc":[3,7,3,1],"tsz":"15.20791328s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","Sys1","STOP"],"etc":[3,7,4,1],"tsz":"15.207927726s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","core","STOP"],"etc":[3,7,5,1],"tsz":"15.207939807s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","otdg-main-wind-up","RWE-34"],"etc":[3,7,6],"msg":"rec-windup-end","tsz":"15.207949304s","tsa":"20250709-001350"}
    {"level":"INFO","tags":["mwt","END"],"etc":[3,8],"tsz":"15.208035499s","tsa":"20250709-001350"}

