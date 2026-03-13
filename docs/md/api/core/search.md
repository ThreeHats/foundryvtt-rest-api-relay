---
title: Search
sidebar_label: Search
---

# Search

## GET /search

Searches Foundry VTT entities using QuickInsert Filters can be a single string for filtering by type ("actor", "item", ext.), or chained together (name:bob,documentType:actor) Available filters: - documentType: type of document ("Actor", "Item", ext) - folder: folder location of the entity (not always defined) - id: unique identifier of the entity - name: name of the entity - package: package identifier the entity belongs to (compendiums minus "Compendium.") - packageName: human-readable package name (readable name of compendium) - subType: sub-type of the entity ("npc", "equipment", ext) - uuid: universal unique identifier - icon: icon HTML for the entity - journalLink: journal link to entity - tagline: same as packageName - formattedMatch: HTML with **applied to matching search parts** - **resultType: constructor name of the QuickInsert result type ("EntitySearchItem". "CompendiumSearchItem", "EmbeddedEntitySearchItem", ext)**

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| filter | string | No | query |  |
| query | string | No | query | Search string |

### Returns

**object** - OK
