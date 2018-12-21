# matrix_creator_morse

This project uses the [MATRIX Creator](https://www.matrix.one/products/creator) to display Morse code with its Everloop LEDs from a given string. It is used in conjunction with a web server to display the IP address of a connected client in the physical world.

To get started, first run the following:

```sh
git clone https://github.com/alexfarra/matrix_creator_morse.git
make
```

To use, run:

```sh
./morse <IPv4 address>
```

For example:

```sh
./morse 192.168.2.2
```

This will make the device light up with RGBW colors corresponding to each section of the IP address (i.e. `RRR.GGG.BBB.WWW`). The LEDs will also blink the IP address back as Morse code.

Note:

libsodium is also needed to compile:

```
sudo apt-get install libsodium-dev
```
