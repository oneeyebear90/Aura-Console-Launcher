# 🎮 Aura Console Launcher

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-Join%20Community-5865F2?style=flat&logo=discord)](https://discord.gg/GA3qsy8vtE)

A lightweight, zero-dependency console proxy designed for the **Aura Minecraft Community**. It allows Xbox, PlayStation, and Nintendo Switch players to seamlessly connect to our crossplay Bedrock server (`bedrock.oneeyebear.net`) via local network (LAN) broadcasting—eliminating the need for third-party mobile apps or network DNS edits.

---

## ✨ Features

- **One-Click Execution:** Hardcoded target address (`bedrock.oneeyebear.net:19132`)—no complex CLI flags or setup required.
- **Cross-Platform:** Pre-compiled binaries available for **Windows**, **macOS** (Intel & Apple Silicon), and **Linux**.
- **100% Safe & Transparent:** Open-source Go code with zero external tracking or hidden dependencies.
- **Lightweight:** Tiny static binary executable (~2 MB memory footprint) that runs silently in your computer's terminal.

---

## 🛠️ How It Works

Console editions of Minecraft (Xbox, PlayStation, Switch) do not provide a native "Add Server" button for custom IPs. However, they continuously scan the local Wi-Fi network for LAN games.
*This may not work for Nintendo Switch*

**Aura Console Launcher** acts as a local bridge on your home network:
1. Listens for UDP broadcast pings on local port `19132`.
2. Responds to your console as a local LAN server named **Aura Minecraft**.
3. Forwards game packets directly between your console and `bedrock.oneeyebear.net`.

---

## 🚀 Quick Start Guide

### Step 1: Download
Head over to the [**Releases Page**](https://github.com/oneeyebear90/Aura-Console-Launcher/releases/latest) and download the file for your operating system:
- **Windows:** `Aura-Launcher-Windows-x64.exe`
- **macOS (Apple Silicon M1/M2/M3/M4):** `Aura-Launcher-Mac-ARM64`
- **macOS (Intel):** `Aura-Launcher-Mac-Intel`
- **Linux:** `Aura-Launcher-Linux-x64`

### Step 2: Launch
1. Ensure your computer is connected to the **same Wi-Fi network** as your console.
2. Launch the downloaded file on your computer:
   - **Windows:** Double-click the `.exe` file. If prompted by Windows Defender, select **Allow Access** for Private Networks.
   - **macOS / Linux:** Open a terminal, grant execution permissions once (`chmod +x Aura-Launcher-Mac-ARM64`), and run `./Aura-Launcher-Mac-ARM64`.

### Step 3: Join on Console
1. Turn on your Xbox, PlayStation, or Switch and launch **Minecraft**.
2. Go to **Play** ➔ **Worlds** tab.
3. Select **Aura Minecraft** to connect! It should be tagged as "Lan world".

> **Note:** Keep the launcher terminal window open on your computer while you are playing on your console.

---

## 💬 Support & Community

Need help getting whitelisted to join the Aura Minecraft server? Follow the Discord link below.
- **Discord:** [Join the Aura Discord Server](https://discord.gg/GA3qsy8vtE)

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.
