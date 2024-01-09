fuse-overlayfs.o: main.c utils.c
	$(CC) $(CFLAGS) -c main.c -o $@ -I. -I/usr/include/fuse3 -I./lib -DPKGLIBEXECDIR=\"fuse-overlayfs\"
	
utils.o: utils.c
	$(CC) $(CFLAGS) -c utils.c -o $@ -I. -I/usr/include/fuse3 -I./lib 
 
direct.o: direct.c
	$(CC) $(CFLAGS) -c direct.c -o $@ -I. -I/usr/include/fuse3 -I./lib

plugin-manager.o: plugin-manager.c
	$(CC) $(CFLAGS) -c plugin-manager.c -o $@ -I. -I/usr/include/fuse3 -I./lib

fuse-overlayfs.a: fuse-overlayfs.o utils.o direct.o plugin-manager.o
	$(AR) rcs $@ $^
	$(RANLIB) $@

build: fuse-overlayfs.a
	echo "build done"

clean:
	rm -f *.o *.a