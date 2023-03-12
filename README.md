# DO NOT USE (if you are sane ofc)
You can change the root directory and if you allow directory listings or not BUT do not use this shit


Still not even good enough for my personal use imma fix it later

Issues;

1. Nil pointers (Happens sometimes and crashes the entire server also idk how to fix it)

2. Some connection issues (mostly during form requests)

3. Goroutine leaks

4. Random errors idk how to fix it but it doesnt crash the entire program so who cares

5. Awful code

6. IndexFirst now works (Plaster fix by adding "/" after cfg.RootDirectory)

I just tested the server with Apache Jmeter (100000 Max users in thread control,hold for 1000 secs,https,GET /,etc) I think its around 12000 Users before my CPU is overloaded and shut off (in task manager i see it uses around 120-170MB of ram?? Ain't no way)

Since its not real world broswer usage. I think the max load is around a 1000 users (maybe?) (nah just ignore this no one will use it)

btw acessing files wtih special chars as name now works
