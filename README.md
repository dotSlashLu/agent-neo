# Agent Neo

An agent running on hosts to control virtual machines remotely via libvirt API.

---

Currently supports:
- volume create, attach, detach, delete
- power operations including suspend, resume
- net interface detach

A few things to notice before using:
- It uses binary protocol for rapid prototyping but might not be very friendly to the implementation of the server
- A demo server written in python is in the server\_demo directory
