# TCam

TCam is an analyser for Tibiantis CAMs.

It implements the following commands:

* `dialogues`: prints the "say"-channel conversations involving the specified name.
* `creature`: prints the timestamps and locations of the specified creature's sightings.
* `location`: prints the timestamps marking the specified location being visited.

Refer to `./tcam help` and `./tcam help [command]` for more information.

## Examples

### Dialogues

```
$ ./tcam --camdir=.../Tibiantis/cam/ dialogues Frodo

# Dialogues for Frodo

## 2024-12-24-11-01-28.cam

### 7m50.453s

* Playa: hi
* Santa Claus: Merry Christmas, little Playa!
* Frodo: Welcome to Frodo's Hut. You heard about the news?
* Playa: rope
* Frodo: Please come back from time to time.
```

### Creature

```
$ ./tcam --camdir=.../Tibiantis/cam/ creature Demon

# Creature Demon

## 2025-10-03-23-11-16.cam

* demon @ (33230, 31643, 14) at 2h19m0.406s

## 2026-01-11-20-28-17.cam

* demon @ (33203, 31640, 15) at 12m17.015s
* demon @ (33204, 31640, 15) at 12m17.265s
* demon @ (33284, 31592, 12) at 16m50.921s
```

### Location

```
$ ./tcam --camdir=.../Tibiantis/cam/ location 33303 31593 14

# Location (x=33303, y=31593, z=14, r=7)

## 2026-01-11-20-28-17.cam

* 18m57s
* 36m0s
* 51m13s
* 1h3m55s
* 1h16m18s
* 1h33m39s
* 1h50m5s
* 2h6m23s
* 2h23m21s
```