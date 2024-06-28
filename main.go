package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	//f := func(w http.ResponseWriter, r *http.Request) {
	//	all, err := io.ReadAll(r.Body)
	//	fmt.Println(string(all), err)
	//	w.Write([]byte("hello world"))
	//	//panic("test panic")
	//}
	//
	//s := http.Server{
	//	Addr:    ":8080",
	//	Handler: http.HandlerFunc(f),
	//	ConnState: func(conn net.Conn, state http.ConnState) {
	//		fmt.Println(state)
	//	},
	//
	//	ConnContext: func(ctx context.Context, c net.Conn) context.Context {
	//		fmt.Println("conn context")
	//		return ctx
	//	},
	//}
	//
	//err := s.ListenAndServe()
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println("hello world")
	//fmt.Println("hello world12222")

	// ./tinydocker run /bin/bash

	switch os.Args[1] {
	case "run":
		fmt.Println("run pid", os.Getpid(), "ppid", os.Getppid())
		readlink, err := os.Readlink("/proc/self/exe")
		if err != nil {
			fmt.Println("readlink()", err.Error())
			return
		}

		fmt.Println(readlink)

		os.Args[1] = "init"
		shCmd := exec.Command(readlink, os.Args[1:]...)
		shCmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		}

		shCmd.Env = os.Environ()
		shCmd.Stdin = os.Stdin
		shCmd.Stdout = os.Stdout
		shCmd.Stderr = os.Stderr

		err = shCmd.Run()
		if err != nil {
			log.Fatal(err.Error())
		}

	case "init":
		cmd := os.Args[2]

		fmt.Println(cmd)

		cwd, err := os.Getwd()
		if err != nil {
			return
		}
		path := cwd + "/ubuntu-base-16.04.6-base-arm64"
		//err := syscall.Chroot("./ubuntu-base-16.04.6-base-arm64")
		//if err != nil {
		//	fmt.Println("chroot()", err.Error())
		//	return
		//}

		//------------------------
		if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
			fmt.Println("mount() /", err.Error())
			return
		}

		if err := syscall.Mount(path, path, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
			fmt.Println("Mount() bind failed,", path, err)
			return
		}

		if err := os.MkdirAll(path+"/.old", 0700); err != nil {
			fmt.Println("mkdir", err)
			return
		}

		if err := syscall.PivotRoot(path, path+"/.old"); err != nil {
			fmt.Println("pivot root ", err)
			return
		}

		if err := syscall.Mount("proc", "/proc", "proc", uintptr(syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV), ""); err != nil {
			fmt.Println("Mount() proc", err.Error())
			return
		}

		//--------------------------------

		if err := syscall.Chdir("/"); err != nil {
			fmt.Println("chdir()", err.Error())
			return
		}

		if err := syscall.Exec(cmd, os.Args[2:], os.Environ()); err != nil {
			fmt.Println("exec()", err)
			return
		}

		fmt.Println("never exec it ")
		return
	}

}
