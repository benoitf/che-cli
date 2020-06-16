# che-cli

Allows to open file in the IDE from a terminal

![screenshot](../assets/che-cli.gif)

## Try out:

1. Open a terminal in Eclipse Che

note: for now it requires this devfile for Che-Theia: https://gist.github.com/benoitf/858289958c13dfa684755d776b0e9c46

2. Grab che binary
```bash
$ wget https://github.com/benoitf/che-cli/releases/download/0.0.1/che
```

3. Set execution flag
```
$ chmod u+x che
```

4. Open a file
```
$ ./che open /projects/<file-to-open>
```
