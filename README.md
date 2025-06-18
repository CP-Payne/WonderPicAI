# 🎨 AI Image Generator Web App - WonderPicAI

A **server-side rendered web application** built in **Go**, allowing users to **generate AI images**, manage a personal gallery, and purchase credits — all within a clean, modular architecture built for scalability and flexibility.


## 🚀 Overview

This portfolio project showcases a modern SSR (server-side rendered) web application using Go and a modular architecture inspired by **Clean Architecture** and **Ports and Adapters (Hexagonal)** principles.

The app is fully functional with:

* 🔐 Authentication (email/password & Google Sign-In)
* 🖼️ Image generation interface powered by AI
* 💳 Stripe integration for purchasing credits
* 🧭 Clean UX with toast notifications, error pages, and a dynamic gallery



## 🧱 Tech Stack

* **Backend:**

  * [Go](https://golang.org/)
  * [Chi](https://github.com/go-chi/chi) – Lightweight, idiomatic HTTP router
  * [Gorm](https://gorm.io/) – ORM for database access
  * [ComfyLite](https://github.com/CP-Payne/ComfyLite) – Modular wrapper for ComfyUI (image generation)

* **Frontend:**

  * [HTMX](https://htmx.org/) – Dynamic frontend interactivity with minimal JS
  * [Templ](https://templ.guide/) – Component-based HTML templating for Go
  * [TailwindCSS](https://tailwindcss.com/) + [daisyUI](https://daisyui.com/) – UI design system and components

* **Architecture:**
  * Ports and Adapters (Hexagonal)


## 📸 Screenshots

### Landing Page
![Landing Page Part 1](web/static/assets/images/screenshots/landing-1.png) 
![Landing Page Part 2](web/static/assets/images/screenshots/landing-part-2.png) 


### Auth pages

![Login Page](web/static/assets/images/screenshots/login-page.png) 
![Signup Page](web/static/assets/images/screenshots/signup-page.png) 


### Image Generation Page

![Image Generation Page](web/static/assets/images/screenshots/gen-page.png) 

### Credits Page

![Credits Page](web/static/assets/images/screenshots/credits-page.png) 

### Stripe Checkout Page
![Checkout Page](web/static/assets/images/screenshots/stripe-checkout.png) 


### Error Page 

![Error Page](web/static/assets/images/screenshots/errorpage.png) 
The error message and heading can easily be modified per error. There is also custom toasts, such as success, info, error and warning, which also has customisable text when called. Furthermore, all input from the user are validated. If validation fails, then appropriate messages will be displayed below the input boxes.


## WonderPicAI C4 Model Diagrams

### 1. System Context Diagram (Level 1)
<p align="center">
  <img src="docs/diagrams/c4/c4-l1-system-context-wonderpicai.drawio.svg" alt="WonderPicAI System Context Diagram" width="800">
</p>

### 2. Container Diagram (Level 2)
<p align="center">
  <img src="docs/diagrams/c4/c4-l2-container-wonderpicai-system.drawio.svg" alt="WonderPicAI Container Diagram" width="900">
</p>

### 3. Component Diagram (Level 3 - SSR Web Application)
<p align="center">
  <img src="docs/diagrams/c4/c4-l3-component-ssr-web-application.drawio.svg" alt="API Application Component Diagram" width="1000">
</p>

## 🧠 Key Features

* 🧩 **Modular Architecture** – Built using interfaces, ensuring that components like AI generation can be swapped without changing the core app.
* 🖼️ **ComfyLite API** – A lightweight, modular API for [ComfyUI](https://www.comfy.org/), designed to be easily replaceable or upgradable.
* 🔒 **Authentication** – Login/Register with password or Google Sign-In (JWT-based)
* 💳 **Credits System** – Stripe integration for purchasing generation credits
* 🖍️ **Gallery** – Each user has a personal gallery to view previously generated images
* ✉️ **Toasts & Error Pages** – User feedback via UI notifications and graceful error handling


## 🔄 Planned Enhancements

While the current app is functional and showcase-ready, I plan to add:

* ✅ **Email Verification**
  * Prevents unintended account access due to unverified ownership
* ⚙️ **Settings Page**
  * Change password, email, and manage account preferences

## 🧰 Project Setup

... PLEASE COMPLETE THIS...



## 🧩 About ComfyLite

[**ComfyLite**](https://github.com/CP-Payne/comfylite) is a lightweight **Go-based REST API wrapper** around [ComfyUI](https://www.comfy.org/), an open-source image generation system.

It was created to:

* Serve as a **decoupled bridge** between the SSR web app and the local image generation engine
* Expose a simple REST API for submitting prompts, image size, and batch settings
* Internally handle **WebSocket communication** with ComfyUI to track generation progress and retrieve results
* Be **model-agnostic** — the web app doesn't care what generation backend is used, as long as it implements the expected interface


Because the web app is built around interfaces (via Ports and Adapters), ComfyLite can be swapped out for a different backend (e.g., a hosted service or another local model) **without changing the core application logic**.

> For implementation details, see the [ComfyLite repository](https://github.com/CP-Payne/ComfyLite)




## 📌 Important Notes

* This project is a **portfolio demonstration**, not a commercial product
* It was built to highlight:

  * My backend engineering skills (Golang, clean architecture)
  * Ability to build full-featured apps with modern frontend UX (HTMX, Tailwind, SSR)
  * Integration of real-world systems (Stripe, Google Auth, AI inference)

## 📫 Contact

If you're a recruiter or fellow developer interested in discussing this project, feel free to reach out:

* LinkedIn: [Charles Payne](https://www.linkedin.com/in/charles-p-payne/)
* Email: [charlpayne1@gmail.com](mailto:charlpayne1@gmail.com)


