---
title: File System
sidebar_label: File System
---

# File System

## GET /download

/download

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query |  |
| format | string | No | query |  |
| path | string | No | query |  |

### Returns

**object** - OK

---

## GET /file-system

/file-system

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query |  |

### Returns

**object** - OK

---

## POST /upload

/upload

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query |  |
| filename | string | No | query |  |
| overwrite | boolean | No | query |  |
| path | string | No | query |  |
| body | string | No | body |  |

### Returns

**object** - Created
