# 🐍 Snake Game in Go (Ebiten)

A classic **Snake Game** built in **Go (Golang)** using the [Ebiten](https://ebiten.org) 2D game library.  
Lightweight, cross-platform, and written entirely in Go — no external assets or engines needed.

---

## 🎮 Features

- Simple, smooth snake movement
- Food spawning and growth
- Increasing speed as you score points
- Collision detection (walls & self)
- Game over + restart system
- Keyboard controls: **Arrow keys** or **WASD**
- Clean and minimalistic graphics

---

## 🧰 Requirements

Make sure you have the following installed:

- **Go** ≥ 1.20
- Linux / macOS / Windows
- For Linux users: you need X11 and OpenGL development libraries

### Linux Dependencies (Ubuntu / Debian / Mint)

```bash
sudo apt-get update
sudo apt-get install -y \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libgl1-mesa-dev \
    xorg-dev \
    libxxf86vm-dev
```
