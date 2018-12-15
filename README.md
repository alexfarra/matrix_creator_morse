# matrix_creator_morse
Project to use the MATRIX Creator to display morse code with the Everloop LEDs from a given string. Used in conjunction with a web server to display the IP address of a connected client in the physical world.

After cloning repo run:

```
make dependencies && make
```

To use:

```
./morse <IPv4 address>
```

For example:

```
./morse 192.168.2.2
```

will make the device light up in RGBW colors corresponding to each section of the IP address. You can think of an IP address as RRR.GGG.BBB.WWW and the LEDs will blink the IP back in morse code.
