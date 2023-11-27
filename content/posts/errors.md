---
title: "Errors"
date: 2023-11-26T10:45:20-05:00
draft: false
---


Over the Thanksgiving week, the atmosphere was delightful with Liz, Brad, and Maggie visiting from out of town along with their families. The children reveled in playing with their cousins. We capped off the week with a grand BBQ at my place on Saturday, inviting Abram and his family as well. Sunday, naturally, became a day of relaxation and recovery after the culinary escapades and late-night festivities.

With a touch of boredom on Sunday, I sought out a productive endeavor. I took a step toward vehicle maintenance by ordering oil filters for Leah's truck, gearing up for an upcoming oil change. Technical tasks beckoned as well—I updated the authorized_keys for vm1 and vm2 on soh.re, facilitating secure shell access from my central server. Engaging Ansible, I meticulously updated all server packages. Subsequently, I delved into [soh.re](https://soh.re), aiming to replicate intermittent crashes by manually terminating some containers with docker kill and then navigating to the website.

```bash
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: 2023/11/26 15:38:39 Accepted connection &{{0xc000086000}}
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: 2023/11/26 15:38:40 dial tcp 127.0.0.1:36897: connect: connection refused
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: 2023/11/26 15:38:40 Connected to localhost &{{0xc000086000}}
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: panic: runtime error: invalid memory address or nil pointer dereference
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x5de9bd]
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: goroutine 3897 [running]:
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: main.forward.func2()
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]:         /home/runner/work/soh-router/soh-router/random.go:33 +0x3d
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]: created by main.forward
Nov 26 15:38:40 vm1.soh.re soh-router[3005682]:         /home/runner/work/soh-router/soh-router/random.go:32 +0x1fc
Nov 26 15:38:40 vm1.soh.re systemd[1]: soh-router.service: Deactivated successfully.
Nov 26 15:38:40 vm1.soh.re systemd[1]: soh-router.service: Consumed 8min 42.538s CPU time.
Nov 26 15:38:41 vm1.soh.re systemd[1]: soh-router.service: Scheduled restart job, restart counter is at 25.
Nov 26 15:38:41 vm1.soh.re systemd[1]: Stopped The Soh-Router service.
Nov 26 15:38:41 vm1.soh.re systemd[1]: soh-router.service: Consumed 8min 42.538s CPU time.
Nov 26 15:38:41 vm1.soh.re systemd[1]: Started The Soh-Router service.
```

Indeed, the failed attempt spurred further exploration. Having recently come across [NilAway](https://www.uber.com/blog/nilaway-practical-nil-panic-detection-for-go/) on Hacker News, I saw potential in its application to the current challenge. Despite its reputation, following the instructions at [https://github.com/uber-go/nilaway](https://github.com/uber-go/nilaway), it yielded no messages specific to this project. Recognizing NilAway as a linter, I expanded my diagnostics to include my standard tools—golint and golangci-lint run. While golint remained silent, golangci-lint surfaced valuable errors, contributing to the debugging process.

```bash
config.go:17:16: Error return value of `yaml.Unmarshal` is not checked (errcheck)
	yaml.Unmarshal(configFile, &v)
	              ^
random.go:30:10: Error return value of `io.Copy` is not checked (errcheck)
		io.Copy(client, conn)
		       ^
random.go:35:10: Error return value of `io.Copy` is not checked (errcheck)
		io.Copy(conn, client)
		       ^
sql3.go:80:15: Error return value of `tx.Rollback` is not checked (errcheck)
			tx.Rollback()
			           ^
sql3.go:111:14: Error return value of `tx.Rollback` is not checked (errcheck)
		tx.Rollback()
		           ^
dockerPool.go:81:9: ineffectual assignment to err (ineffassign)
		rows, err := db.Query("SELECT name FROM docker_pool;")
		      ^
config.go:5:2: SA1019: "io/ioutil" has been deprecated since Go 1.19: As of Go 1.16, the same functionality is now provided by package [io] or package [os], and those implementations should be preferred in new code. See the specific function documentation for details. (staticcheck)
	"io/ioutil"
	^
```

Even after addressing the issues identified by golangci-lint, the persistent recurrence of the same error prompted a deeper investigation. My scrutiny revealed an oversight in error handling—specifically, the neglect to check the return values of client.Close and conn.Close within deferred functions. Seeking guidance, I reached out to [ChatGPT](https://chat.openai.com/share/e07fdda1-206f-4080-910f-0e4b34de5a76) for insights on the optimal approach for error checking within a defer statement. The response included an ingenious code snippet, suggesting the deferral of an entire new function to handle errors. Implementation of this new defer code unveiled previously overlooked errors, shedding light on the intricacies of the problem.

```bash
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: 2023/11/26 16:15:23 Accepted connection &{{0xc000128700}}
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: 2023/11/26 16:15:23 dial tcp 127.0.0.1:33029: connect: connection refused
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: 2023/11/26 16:15:23 Connected to localhost &{{0xc000128700}}
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: panic: runtime error: invalid memory address or nil pointer dereference
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         panic: runtime error: invalid memory address or nil pointer dereference
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: [signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x5d245a]
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: goroutine 184 [running]:
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: main.forward.func2.1()
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /home/jmainguy/Github/soh-router/random.go:45 +0x1a
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: panic({0x6d7f60?, 0x8c2320?})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/runtime/panic.go:914 +0x21f
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: io.copyBuffer({0x754320, 0xc0002f6720}, {0x0, 0x0}, {0x0, 0x0, 0x0})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/io/io.go:430 +0x18b
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: io.Copy(...)
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/io/io.go:389
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: net.genericReadFrom({0x7543c0?, 0xc0001109c0?}, {0x0, 0x0})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/net/net.go:671 +0x66
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: net.(*TCPConn).readFrom(0xc0001109c0, {0x0, 0x0})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/net/tcpsock_posix.go:54 +0x6b
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: net.(*TCPConn).ReadFrom(0xc0001109c0, {0x0?, 0x0?})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/net/tcpsock.go:130 +0x30
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: io.copyBuffer({0x7543c0, 0xc0001109c0}, {0x0, 0x0}, {0x0, 0x0, 0x0})
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/io/io.go:416 +0x147
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: io.Copy(...)
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /usr/lib/golang/src/io/io.go:389
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: main.forward.func2()
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /home/jmainguy/Github/soh-router/random.go:56 +0x105
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]: created by main.forward in goroutine 180
Nov 26 16:15:23 vm1.soh.re soh-router[3221059]:         /home/jmainguy/Github/soh-router/random.go:43 +0x231
```

Delving deeper into the code with the insights gleaned from the updated logs, a glaring revelation emerged—we were attempting to send traffic to a target that did not exist. This pivotal realization shed light on the root cause of the persistent error and paved the way for a more targeted and effective resolution.

```go
target := pullDockerFromPool(db)
client, err := net.Dial("tcp", target)
if err != nil {
	check(err)
}
log.Printf("Connected to localhost %v\n", conn)
```


Recognizing a critical flaw in our error-checking approach, where the check merely printed a log statement, we consistently received a misleading "Connected to localhost %v" message even when encountering connection errors. To rectify this, I introduced a significant change by encapsulating the entire process within a for loop. This alteration ensured persistent attempts until a successful connection was established. As a result, the server demonstrated increased resilience, no longer succumbing to crashes under similar circumstances.

```go
var client net.Conn
var err error
var attempts int
for {
	attempts++
	// Pull a container from the pool, and connect to it
	target := pullDockerFromPool(db)
	client, err = net.Dial("tcp", target)
	if err != nil {
		// We were unable to connect to the container, try a different one
		log.Printf("Attempt #%d for target has failed\n", attempts)
		log.Println(err)
	} else {
		// End for loop, we got a good target
		break
	}
}
log.Printf("Connected to localhost %v\n", conn)
```

The observed anomaly in the number of connection attempts—25 instead of the anticipated 10—raised concerns, especially considering the minimal traffic on the site, primarily driven by personal use. A closer examination uncovered a staggering 259 containers in the database pool designated to receive traffic. This stark deviation from expectations signaled a potential malfunction in the reaper.

Suspecting a discrepancy in string comparisons, I initiated an investigation into the code responsible for managing the reaper to pinpoint the root cause of this unexpected behavior.

```go
func reap(db *sql.DB, name string) {
	// If container does not exist, remove from pool)
	running, err := exec.Command("docker", "inspect", "--format='{{.State.Running}}'", name).Output()
	check(err)
	isRunning := string(running)
	if err != nil {
		DelName(db, name)
		log.Printf("Reaped %v", name)
	}
	if isRunning == "false\n" {
		DelName(db, name)
		log.Printf("Reaped %v", name)
	}
}
```

After adding some more logging I saw that the string being returned for false was 

```bash
Nov 26 16:54:57 vm1.soh.re soh-router[3394994]: 2023/11/26 16:54:57 Did not reap 93061D4B16339792BAAF because isRunning returned ''false'
Nov 26 16:54:57 vm1.soh.re soh-router[3394994]: '
```

This revelation unveiled that the string returned by Podman was 'false'\n. Rectifying the code to accommodate this format corrected the reaper's functionality. Implementation of the updated code on both VMs revealed a significant improvement, with VM2 no longer failing to reap numerous connections as observed in VM1.

Remarkably, the system demonstrated resilience by functioning even when unable to write to the database on disk. This characteristic was highlighted during the installation of the new code. To ensure seamless future deployments, I created a configuration file, /etc/soh-router/config.yaml, and established the database directory at /opt/soh-router. Consequently, the database began to be written to disk as expected. Additional code for these configurations was incorporated into the .gorelease.yml for streamlined assistance in future releases.

[Link to pull request on GitHub](https://github.com/Jmainguy/soh-router/pull/13)
