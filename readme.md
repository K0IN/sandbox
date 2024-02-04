
# Sandbox

A small go program to create a sandbox of your current system.
Try out commands without the fear of breaking your system.

## Install

> curl -sSL https://raw.githubusercontent.com/K0IN/sandbox/main/install.sh | sh

which will install the `sandbox` command to `/usr/local/bin`.

> [!WARNING]  
> This is a setuid binary, this is needed to create a overlayfs mount of /

## Usage

> sandbox --help

## Examples

Run a command inside a sandbox:

> sandbox try "<your command>"

Run a interactive shell inside a sandbox:

> sandbox try

Persist the sandbox:

> sandbox try --persist "<your command>"

Example on what you can do with persist:

```shell
k0in@K0IN-PC /m/e/s/sandbox (main)> sandbox try --persist bash
root@K0IN-PC:/mnt/e/src/sandbox# cat > a.txt <<EOF
Hello World
EOF
root@K0IN-PC:/mnt/e/src/sandbox# cat a.txt 
Hello World
root@K0IN-PC:/mnt/e/src/sandbox# exit
exit
Sandbox created at /root/.sandboxes/iQXpuc
k0in@K0IN-PC /m/e/s/sandbox (main)> sandbox status iQXpuc
not staged: mnt/e/src/sandbox/a.txt
not staged: root/.bash_history
k0in@K0IN-PC /m/e/s/sandbox (main)> sandbox diff iQXpuc
--- /mnt/e/src/sandbox/a.txt
+++ /root/.sandboxes/iQXpuc/upper/mnt/e/src/sandbox/a.txt
@@ -1 +1 @@
+Hello World

--- /root/.bash_history
+++ /root/.sandboxes/iQXpuc/upper/root/.bash_history
@@ -12,3 +12,9 @@
 cd home/
 ll
 exit
+cat > a.txt <<EOF
+Hello World
+EOF
+
+cat a.txt 
+exit

k0in@K0IN-PC /m/e/s/sandbox (main)> sandbox add iQXpuc mnt/e/src/sandbox/a.txt
Staged mnt/e/src/sandbox/a.txt
k0in@K0IN-PC /m/e/s/sandbox (main)> sandbox status iQXpuc
staged: mnt/e/src/sandbox/a.txt
not staged: root/.bash_history
k0in@K0IN-PC /m/e/s/sandbox (main) [1]> sandbox commit iQXpuc
Are you sure you want to commit? [y/N]: y
k0in@K0IN-PC /m/e/s/sandbox (main)> ls
LICENSE  VP-XHr  a.txt  cli  go.mod  go.sum  helper  install.sh  main.go  readme.md  sandbox
k0in@K0IN-PC /m/e/s/sandbox (main)> cat a.txt 
Hello World
```



## Credits

This is highly inspired by [try](https://github.com/binpash/try/blob/main/try) - give it a try!
