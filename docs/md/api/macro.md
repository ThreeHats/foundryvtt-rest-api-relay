---
tag: macro
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Macro

## GET /macros

Get all macros

Retrieves a list of all macros available in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of macros with details

### Try It Out

<ApiTester
  method="GET"
  path="/macros"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/macros';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'GET',
  headers: {
    'x-api-key': 'your-api-key-here'
  }
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X GET 'http://localhost:3010/macros?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/macros'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.get(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    }
)
data = response.json()
print(data)
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import axios from 'axios';

(async () => {
  const baseUrl = 'http://localhost:3010';
  const path = '/macros';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'get',
    headers: {
      'x-api-key': 'your-api-key-here'
    },
    url
  });
  const data = response.data;
  console.log(data);
})();
```

</TabItem>
<TabItem value="emojicode" label="Emojicode">

```emojicode
📦 sockets 🏠

💭 Emojicode HTTP Client
💭 Compile: emojicodec example.🍇 -o example
💭 Run: ./example

🏁 🍇
  💭 Connection settings
  🔤localhost🔤 ➡️ host
  3010 ➡️ port
  🔤/macros🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /macros🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

  💭 Connect and send
  🍺 🆕📞 host port❗ ➡️ socket
  🍺 💬 socket 📇 request❗❗
  
  💭 Read and print response
  🍺 👂 socket 4096❗ ➡️ data
  😀 🍺 🔡 data❗❗
  
  💭 Close socket
  🚪 socket❗
🍉
```

</TabItem>
</Tabs>

#### Response

**Status:** 200

```json
{
  "type": "macros-result",
  "requestId": "macros_1776657989837",
  "macros": [
    {
      "uuid": "Macro.a8GD8bNBNJQeVtd1",
      "id": "a8GD8bNBNJQeVtd1",
      "name": "s2s",
      "type": "script",
      "author": "Gamemaster",
      "command": "// ============================================================================\n// Server-to-Server Transfer Macro (Phase 1)\n// Transfers entities between Foundry VTT worlds via the REST API relay server.\n// Source = this server (local Foundry APIs). Target = remote server (REST API).\n// Requires: foundry-rest-api module installed and configured on both worlds.\n// ============================================================================\n(async () => {\n  \"use strict\";\n\n  // ========================= SECTION 1: CONFIG & HELPERS =========================\n\n  const MODULE_ID = \"foundry-rest-api\";\n  const ENTITY_TYPES = [\"Actor\", \"Item\", \"Scene\", \"JournalEntry\", \"RollTable\", \"Cards\", \"Macro\", \"Playlist\"];\n  const FILE_EXTENSIONS = /\\.(png|jpg|jpeg|gif|webp|svg|avif|mp3|ogg|wav|flac|m4a|webm|mp4|pdf)$/i;\n  const SKIP_PATH_PREFIXES = [\"systems/\", \"modules/\", \"icons/\", \"ui/\"];\n\n  // Map entity type names to game collections\n  const COLLECTION_MAP = {\n    Actor: game.actors,\n    Item: game.items,\n    Scene: game.scenes,\n    JournalEntry: game.journal,\n    RollTable: game.tables,\n    Cards: game.cards,\n    Macro: game.macros,\n    Playlist: game.playlists,\n  };\n\n  // Validate module is active\n  const mod = game.modules.get(MODULE_ID);\n  if (!mod?.active) {\n    ui.notifications.error(\"Foundry REST API module is not installed or active.\");\n    return;\n  }\n\n  // Read settings\n  const wsUrl = game.settings.get(MODULE_ID, \"wsRelayUrl\");\n  const apiKey = game.settings.get(MODULE_ID, \"apiKey\");\n\n  if (!wsUrl || !apiKey) {\n    ui.notifications.error(\"REST API module is not configured. Please set the relay URL and API key.\");\n    return;\n  }\n\n  // Derive HTTP URL from WebSocket URL\n  const relayUrl = wsUrl.replace(\"wss://\", \"https://\").replace(\"ws://\", \"http://\").replace(/\\/+$/, \"\");\n\n  // ========================= SECTION 2: API CLIENT (target server only) =========================\n\n  async function relayFetch(endpoint, options = {}) {\n    const url = new URL(endpoint, relayUrl);\n    if (options.params) {\n      for (const [k, v] of Object.entries(options.params)) {\n        if (v !== undefined && v !== null) url.searchParams.set(k, String(v));\n      }\n    }\n\n    const headers = { \"x-api-key\": apiKey };\n    if (options.body) headers[\"Content-Type\"] = \"application/json\";\n\n    const resp = await fetch(url.toString(), {\n      method: options.method || \"GET\",\n      headers,\n      body: options.body ? JSON.stringify(options.body) : undefined,\n    });\n\n    if (!resp.ok) {\n      const text = await resp.text().catch(() => \"\");\n      throw new Error(`API ${resp.status}: ${text || resp.statusText}`);\n    }\n\n    return resp.json();\n  }\n\n  // Remote API — only used for TARGET server operations\n  const remote = {\n    getClients: () => relayFetch(\"/clients\"),\n\n    createEntity: (entityType, data, clientId, folder, keepId = true) =>\n      relayFetch(\"/create\", {\n        method: \"POST\",\n        params: { clientId },\n        body: { entityType, data, folder: folder || null, keepId },\n      }),\n\n    uploadFile: (path, filename, fileData, clientId) =>\n      relayFetch(\"/upload\", {\n        method: \"POST\",\n        params: { clientId },\n        body: { path, filename, fileData, source: \"data\", overwrite: true },\n      }),\n\n    createFolder: (name, folderType, clientId, parentFolderId) =>\n      relayFetch(\"/create-folder\", {\n        method: \"POST\",\n        params: { clientId, name, folderType, parentFolderId },\n      }),\n\n    getStructure: (clientId, types, includeEntityData = false, recursive = true) =>\n      relayFetch(\"/structure\", {\n        params: { clientId, types: types.join(\",\"), includeEntityData, recursive, recursiveDepth: 10 },\n      }),\n\n    getUsers: (clientId) => relayFetch(\"/users\", { params: { clientId } }),\n\n    createUser: (name, role, password, clientId) =>\n      relayFetch(\"/user\", {\n        method: \"POST\",\n        params: { clientId },\n        body: { name, role: role ?? 1, password: password || undefined },\n      }),\n  };\n\n  // ========================= SECTION 3: LOCAL SOURCE HELPERS =========================\n\n  // Get folder chain (array of folder names from root to the entity's folder)\n  function getFolderChain(entity) {\n    const chain = [];\n    let folder = entity.folder;\n    while (folder) {\n      chain.unshift(folder.name);\n      folder = folder.folder;\n    }\n    return chain;\n  }\n\n  // Serialize an entity for transfer using Foundry's toObject\n  function serializeEntity(entity) {\n    return entity.toObject(true);\n  }\n\n  // Download a local file as a data URL\n  async function downloadLocalFile(filePath) {\n    const url = filePath.startsWith(\"http\") ? filePath : foundry.utils.getRoute(filePath);\n    const resp = await fetch(url);\n    if (!resp.ok) throw new Error(`${resp.status} ${resp.statusText}`);\n    const blob = await resp.blob();\n    return new Promise((resolve, reject) => {\n      const reader = new FileReader();\n      reader.onload = () => resolve(reader.result);\n      reader.onerror = reject;\n      reader.readAsDataURL(blob);\n    });\n  }\n\n  // Find all tokens for an actor across all scenes\n  function findActorTokens(actorId) {\n    const tokens = [];\n    for (const scene of game.scenes) {\n      for (const token of scene.tokens) {\n        if (token.actorId === actorId) {\n          tokens.push({ scene, token });\n        }\n      }\n    }\n    return tokens;\n  }\n\n  // Delete all tokens for an actor across all scenes\n  async function deleteActorTokens(actorId, log) {\n    const tokenEntries = findActorTokens(actorId);\n    if (tokenEntries.length === 0) return 0;\n\n    log(`  Deleting ${tokenEntries.length} token(s) across ${new Set(tokenEntries.map(t => t.scene.id)).size} scene(s)`);\n\n    // Group by scene for batch deletion\n    const byScene = new Map();\n    for (const { scene, token } of tokenEntries) {\n      if (!byScene.has(scene.id)) byScene.set(scene.id, { scene, tokenIds: [] });\n      byScene.get(scene.id).tokenIds.push(token.id);\n    }\n\n    let deleted = 0;\n    for (const { scene, tokenIds } of byScene.values()) {\n      try {\n        await scene.deleteEmbeddedDocuments(\"Token\", tokenIds);\n        deleted += tokenIds.length;\n      } catch (err) {\n        log(`  Warning: Failed to delete tokens in scene \"${scene.name}\": ${err.message}`);\n      }\n    }\n\n    return deleted;\n  }\n\n  // ========================= SECTION 4: FILE TRANSFER ENGINE =========================\n\n  function extractFilePaths(data) {\n    const paths = new Set();\n\n    function walk(obj) {\n      if (!obj || typeof obj !== \"object\") {\n        if (typeof obj === \"string\" && FILE_EXTENSIONS.test(obj)) {\n          if (obj.startsWith(\"http://\") || obj.startsWith(\"https://\")) return;\n          if (SKIP_PATH_PREFIXES.some((p) => obj.startsWith(p))) return;\n          paths.add(obj);\n        }\n        return;\n      }\n      if (Array.isArray(obj)) {\n        for (const item of obj) walk(item);\n      } else {\n        for (const val of Object.values(obj)) walk(val);\n      }\n    }\n\n    walk(data);\n    return [...paths];\n  }\n\n  async function transferFiles(paths, targetClientId, log) {\n    const pathMap = new Map();\n    let success = 0;\n    let failed = 0;\n\n    for (let i = 0; i < paths.length; i++) {\n      const filePath = paths[i];\n      const parts = filePath.split(\"/\");\n      const filename = parts.pop();\n      const dir = parts.join(\"/\") || \".\";\n\n      log(`  File ${i + 1}/${paths.length}: ${filePath}`);\n\n      try {\n        // Download from local Foundry server\n        const fileData = await downloadLocalFile(filePath);\n\n        // Upload to remote target\n        await remote.uploadFile(dir, filename, fileData, targetClientId);\n        pathMap.set(filePath, filePath);\n        success++;\n      } catch (err) {\n        log(`    Warning: Failed to transfer: ${err.message}`);\n        failed++;\n      }\n    }\n\n    log(`  Files: ${success} transferred, ${failed} failed`);\n    return pathMap;\n  }\n\n  function remapFilePaths(data, pathMap) {\n    if (!pathMap.size) return data;\n\n    function walk(obj) {\n      if (!obj || typeof obj !== \"object\") {\n        if (typeof obj === \"string\" && pathMap.has(obj)) return pathMap.get(obj);\n        return obj;\n      }\n      if (Array.isArray(obj)) return obj.map(walk);\n      const result = {};\n      for (const [k, v] of Object.entries(obj)) result[k] = walk(v);\n      return result;\n    }\n\n    return walk(data);\n  }\n\n  // ========================= SECTION 5: ENTITY TRANSFER ENGINE =========================\n\n  // Cache for created folders on target\n  const folderCache = new Map();\n\n  async function ensureFolderChain(folderChain, folderType, targetClientId, log) {\n    if (!folderChain.length) return null;\n\n    let parentFolderId = null;\n\n    for (const folderName of folderChain) {\n      const cacheKey = `${folderName}|${folderType}|${parentFolderId || \"root\"}`;\n\n      if (folderCache.has(cacheKey)) {\n        parentFolderId = folderCache.get(cacheKey);\n        continue;\n      }\n\n      try {\n        log(`  Creating folder: ${folderName}`);\n        const result = await remote.createFolder(folderName, folderType, targetClientId, parentFolderId);\n        parentFolderId = result.data?.id || result.id;\n        folderCache.set(cacheKey, parentFolderId);\n      } catch (err) {\n        // Folder might already exist — try to find it\n        try {\n          const structure = await remote.getStructure(targetClientId, [folderType], false, true);\n          const folderId = findFolderInStructure(structure, folderName, parentFolderId);\n          if (folderId) {\n            parentFolderId = folderId;\n            folderCache.set(cacheKey, parentFolderId);\n          } else {\n            log(`  Warning: Could not create or find folder \"${folderName}\": ${err.message}`);\n            return parentFolderId;\n          }\n        } catch {\n          log(`  Warning: Could not create folder \"${folderName}\": ${err.message}`);\n          return parentFolderId;\n        }\n      }\n    }\n\n    return parentFolderId;\n  }\n\n  function findFolderInStructure(structure, name, parentId) {\n    const folders = structure?.data?.folders || structure?.folders || {};\n\n    function search(obj) {\n      for (const [key, val] of Object.entries(obj)) {\n        if (!val || typeof val !== \"object\") continue;\n        if (key === name && val.id) {\n          if (parentId && val.parentFolder !== parentId) continue;\n          return val.id;\n        }\n        const nested = search(val);\n        if (nested) return nested;\n      }\n      return null;\n    }\n\n    return search(folders);\n  }\n\n  async function transferEntity(entity, entityType, targetClientId, options, log) {\n    const { deleteFromSource, deleteTokens, transferFilesOpt, ownershipRemap } = options;\n    const entityName = entity.name || entity.id;\n\n    log(`Transferring ${entityType}: ${entityName}`);\n\n    try {\n      // Serialize from local Foundry\n      const entityData = serializeEntity(entity);\n\n      // Transfer files\n      let pathMap = new Map();\n      if (transferFilesOpt) {\n        const filePaths = extractFilePaths(entityData);\n        if (filePaths.length > 0) {\n          log(`  Found ${filePaths.length} file(s) to transfer`);\n          pathMap = await transferFiles(filePaths, targetClientId, log);\n        }\n      }\n\n      // Clone and remap file paths\n      let createData = JSON.parse(JSON.stringify(entityData));\n      createData = remapFilePaths(createData, pathMap);\n\n      // Remap ownership if needed\n      if (ownershipRemap && createData.ownership) {\n        const { sourceUserId, targetUserId } = ownershipRemap;\n        if (createData.ownership[sourceUserId] !== undefined) {\n          const permLevel = createData.ownership[sourceUserId];\n          delete createData.ownership[sourceUserId];\n          createData.ownership[targetUserId] = permLevel;\n        }\n      }\n\n      // Recreate folder hierarchy on target\n      let targetFolderId = null;\n      const folderChain = getFolderChain(entity);\n      if (folderChain.length > 0) {\n        targetFolderId = await ensureFolderChain(folderChain, entityType, targetClientId, log);\n      }\n\n      // Remove folder from data (we pass it separately to the API)\n      delete createData.folder;\n\n      // Create entity on target (preserving _id)\n      const result = await remote.createEntity(entityType, createData, targetClientId, targetFolderId);\n\n      const targetUuid = result.uuid || `${entityType}.${createData._id}`;\n      log(`  Created on target: ${targetUuid}`);\n\n      // Delete from source if requested\n      if (deleteFromSource) {\n        try {\n          // Delete associated tokens first if this is an actor\n          if (deleteTokens && entityType === \"Actor\") {\n            const count = await deleteActorTokens(entity.id, log);\n            if (count > 0) log(`  Deleted ${count} token(s) from source scenes`);\n          }\n\n          await entity.delete();\n          log(`  Deleted from source`);\n        } catch (err) {\n          log(`  Warning: Failed to delete from source: ${err.message}`);\n        }\n      }\n\n      return { success: true, entityName, targetUuid };\n    } catch (err) {\n      log(`  ERROR: ${err.message}`);\n      return { success: false, entityName, error: err.message };\n    }\n  }\n\n  // ========================= SECTION 6: USER TRANSFER ENGINE =========================\n\n  async function transferPlayer(user, targetClientId, options, log) {\n    const { createAccount, password, role, transferFilesOpt, deleteFromSource, deleteTokens, deleteUser } = options;\n\n    let targetUserId = null;\n    let ownershipRemap = null;\n\n    if (createAccount) {\n      log(`Creating user account: ${user.name}`);\n\n      // Check if user already exists on target\n      const targetUsersResp = await remote.getUsers(targetClientId);\n      const targetUsers = targetUsersResp.data || targetUsersResp;\n      const existingUser = (Array.isArray(targetUsers) ? targetUsers : []).find(\n        (u) => u.name.toLowerCase() === user.name.toLowerCase()\n      );\n\n      if (existingUser) {\n        targetUserId = existingUser.id;\n        log(`  User already exists on target (ID: ${targetUserId})`);\n      } else {\n        const created = await remote.createUser(\n          user.name,\n          role ?? user.role,\n          password || undefined,\n          targetClientId\n        );\n        targetUserId = created.data?.id || created.id;\n        log(`  Created user on target (ID: ${targetUserId})`);\n      }\n\n      ownershipRemap = { sourceUserId: user.id, targetUserId };\n    }\n\n    // Find all entities owned by this user (locally)\n    log(\"Finding owned entities...\");\n    const ownedEntities = [];\n    for (const entityType of ENTITY_TYPES) {\n      const collection = COLLECTION_MAP[entityType];\n      if (!collection) continue;\n      for (const entity of collection) {\n        const ownership = entity.ownership || {};\n        if (ownership[user.id] >= 3) {\n          ownedEntities.push({ entity, entityType });\n        }\n      }\n    }\n\n    log(`Found ${ownedEntities.length} owned entity(ies)`);\n\n    // Transfer each entity\n    const results = [];\n    for (const { entity, entityType } of ownedEntities) {\n      const result = await transferEntity(entity, entityType, targetClientId, {\n        deleteFromSource,\n        deleteTokens,\n        transferFilesOpt,\n        ownershipRemap,\n      }, log);\n      results.push(result);\n    }\n\n    // Delete user from source if requested\n    if (deleteUser && deleteFromSource) {\n      try {\n        await user.delete();\n        log(`Deleted user \"${user.name}\" from source`);\n      } catch (err) {\n        log(`Warning: Failed to delete user from source: ${err.message}`);\n      }\n    }\n\n    return results;\n  }\n\n  // ========================= SECTION 7: UI DIALOGS =========================\n\n  // Fetch connected clients to find target servers\n  let clients;\n  try {\n    const resp = await remote.getClients();\n    clients = resp.clients || [];\n  } catch (err) {\n    ui.notifications.error(`Failed to connect to relay server: ${err.message}`);\n    return;\n  }\n\n  if (clients.length < 2) {\n    ui.notifications.warn(\"Need at least 2 connected servers. Ensure both worlds have the REST API module configured with the same API key.\");\n    return;\n  }\n\n  // Source is always the current world\n  const currentWorldId = game.world.id;\n  const currentClient = clients.find((c) => c.worldId === currentWorldId);\n  if (!currentClient) {\n    ui.notifications.error(\"Could not identify this server among connected clients. Check that the REST API module is connected.\");\n    return;\n  }\n  const sourceName = currentClient.customName || currentClient.worldTitle || currentClient.worldId;\n\n  // Target candidates are all other connected servers\n  const targetClients = clients.filter((c) => c.id !== currentClient.id);\n\n  // Shared CSS for all transfer dialogs\n  const S2S_CSS = `\n    <style>\n      .s2s .form-group { margin: 4px 0; }\n      .s2s .form-group label { margin-bottom: 2px; }\n      .s2s .s2s-row { display: flex; gap: 8px; }\n      .s2s .s2s-row > .form-group { flex: 1; }\n      .s2s .s2s-header { display: flex; align-items: center; gap: 6px; margin-bottom: 8px; padding: 6px 8px; background: rgba(0,0,0,0.1); border-radius: 4px; font-size: 13px; }\n      .s2s .s2s-header .fas { opacity: 0.5; }\n      .s2s .s2s-checks { display: flex; flex-wrap: wrap; gap: 4px 16px; margin-top: 6px; padding-top: 6px; border-top: 1px solid rgba(255,255,255,0.1); }\n      .s2s .s2s-checks label { display: flex; align-items: center; gap: 4px; white-space: nowrap; font-size: 12px; }\n      .s2s .s2s-entity-list { max-height: 250px; overflow-y: auto; border: 1px solid rgba(255,255,255,0.15); border-radius: 3px; padding: 4px; }\n      .s2s .s2s-entity-list label { display: block; padding: 1px 2px; font-size: 12px; }\n      .s2s .s2s-entity-list label:hover { background: rgba(255,255,255,0.05); }\n      .s2s .s2s-select-all { display: block; padding: 2px; margin-bottom: 2px; border-bottom: 1px solid rgba(255,255,255,0.1); font-weight: bold; font-size: 12px; }\n      .s2s .s2s-log { height: 200px; overflow-y: auto; background: #0a0a0a; color: #0f0; padding: 8px; font-size: 11px; white-space: pre-wrap; border-radius: 3px; border: 1px solid rgba(255,255,255,0.1); }\n      .s2s .s2s-status { font-size: 12px; margin-bottom: 4px; opacity: 0.8; }\n    </style>\n  `;\n\n  // --- Server Selection Dialog ---\n  function showServerSelectDialog() {\n    const targetOptions = targetClients\n      .map((c) => `<option value=\"${c.id}\">${c.customName || c.worldTitle || c.worldId} (${c.systemTitle || c.systemId})</option>`)\n      .join(\"\");\n\n    const content = `\n      ${S2S_CSS}\n      <div class=\"s2s\">\n        <div class=\"s2s-row\">\n          <div class=\"form-group\">\n            <label>Source</label>\n            <input type=\"text\" value=\"${sourceName}\" disabled />\n          </div>\n          <div class=\"form-group\">\n            <label>Target</label>\n            <select name=\"target\">${targetOptions}</select>\n          </div>\n        </div>\n        <div class=\"form-group\">\n          <label>Transfer Mode</label>\n          <select name=\"mode\">\n            <option value=\"entities\">Select Entities</option>\n            <option value=\"player\">Transfer Player</option>\n          </select>\n        </div>\n      </div>\n    `;\n\n    new Dialog({\n      title: \"Server-to-Server Transfer\",\n      content,\n      buttons: {\n        next: {\n          icon: '<i class=\"fas fa-arrow-right\"></i>',\n          label: \"Next\",\n          callback: (html) => {\n            const targetId = html.find('[name=\"target\"]').val();\n            const mode = html.find('[name=\"mode\"]').val();\n\n            if (mode === \"entities\") showEntitySelectDialog(targetId);\n            else showPlayerTransferDialog(targetId);\n          },\n        },\n        cancel: { icon: '<i class=\"fas fa-times\"></i>', label: \"Cancel\" },\n      },\n      default: \"next\",\n    }, { width: 420 }).render(true);\n  }\n\n  // --- Entity Selection Dialog ---\n  function showEntitySelectDialog(targetId) {\n    const targetClient = clients.find((c) => c.id === targetId);\n    const targetName = targetClient?.customName || targetClient?.worldTitle || targetId;\n\n    const typeOptions = ENTITY_TYPES.map((t) => `<option value=\"${t}\">${t}</option>`).join(\"\");\n\n    const content = `\n      ${S2S_CSS}\n      <div class=\"s2s\">\n        <div class=\"s2s-header\">\n          <span>${sourceName}</span>\n          <i class=\"fas fa-arrow-right\"></i>\n          <span>${targetName}</span>\n        </div>\n        <div class=\"form-group\">\n          <label>Entity Type</label>\n          <select name=\"entityType\" id=\"s2s-entity-type\">${typeOptions}</select>\n        </div>\n        <div class=\"form-group\">\n          <div id=\"s2s-entity-list\" class=\"s2s-entity-list\">\n            <p><em>Loading...</em></p>\n          </div>\n        </div>\n        <div class=\"s2s-checks\">\n          <label><input type=\"checkbox\" name=\"transferFiles\" checked /> Transfer files</label>\n          <label><input type=\"checkbox\" name=\"deleteSource\" /> Delete from source</label>\n          <label id=\"s2s-delete-tokens-group\" style=\"display:none;\"><input type=\"checkbox\" name=\"deleteTokens\" checked /> Delete tokens</label>\n        </div>\n      </div>\n    `;\n\n    // Map of id -> local entity reference\n    let entityRefMap = new Map();\n\n    new Dialog({\n      title: \"Select Entities to Transfer\",\n      content,\n      buttons: {\n        transfer: {\n          icon: '<i class=\"fas fa-exchange-alt\"></i>',\n          label: \"Transfer\",\n          callback: async (html) => {\n            const checked = html.find('.s2s-entity-cb:checked');\n            const selectedIds = [];\n            checked.each(function () { selectedIds.push(this.value); });\n\n            if (selectedIds.length === 0) {\n              ui.notifications.warn(\"No entities selected.\");\n              return;\n            }\n\n            const entityType = html.find('[name=\"entityType\"]').val();\n            const transferFilesOpt = html.find('[name=\"transferFiles\"]').is(\":checked\");\n            const deleteFromSource = html.find('[name=\"deleteSource\"]').is(\":checked\");\n            const deleteTokens = html.find('[name=\"deleteTokens\"]').is(\":checked\");\n\n            // Resolve local entity references\n            const entitiesToTransfer = [];\n            for (const id of selectedIds) {\n              const entity = entityRefMap.get(id);\n              if (entity) entitiesToTransfer.push({ entity, entityType });\n              else ui.notifications.warn(`Entity ${id} not found locally.`);\n            }\n\n            showProgressDialog(entitiesToTransfer, targetId, { transferFilesOpt, deleteFromSource, deleteTokens });\n          },\n        },\n        cancel: { icon: '<i class=\"fas fa-times\"></i>', label: \"Cancel\" },\n      },\n      default: \"transfer\",\n      render: (html) => {\n        const typeSelect = html.find(\"#s2s-entity-type\");\n        const deleteSourceCb = html.find('[name=\"deleteSource\"]');\n        const deleteTokensGroup = html.find(\"#s2s-delete-tokens-group\");\n\n        // Show/hide \"delete tokens\" based on entity type and delete checkbox\n        const updateTokensVisibility = () => {\n          const isActor = typeSelect.val() === \"Actor\";\n          const isDelete = deleteSourceCb.is(\":checked\");\n          deleteTokensGroup.toggle(isActor && isDelete);\n        };\n        typeSelect.on(\"change\", updateTokensVisibility);\n        deleteSourceCb.on(\"change\", updateTokensVisibility);\n\n        const loadEntities = () => {\n          const entityType = typeSelect.val();\n          const listDiv = html.find(\"#s2s-entity-list\");\n          entityRefMap.clear();\n\n          const collection = COLLECTION_MAP[entityType];\n          if (!collection || collection.size === 0) {\n            listDiv.html(\"<p><em>No entities found</em></p>\");\n            updateTokensVisibility();\n            return;\n          }\n\n          // Build list with folder paths\n          const entries = [];\n          for (const entity of collection) {\n            const folderChain = getFolderChain(entity);\n            const folderPath = folderChain.join(\"/\");\n            entries.push({ id: entity.id, name: entity.name, folderPath, entity });\n            entityRefMap.set(entity.id, entity);\n          }\n\n          // Sort by folder path then name\n          entries.sort((a, b) => {\n            const pathCmp = a.folderPath.localeCompare(b.folderPath);\n            if (pathCmp !== 0) return pathCmp;\n            return a.name.localeCompare(b.name);\n          });\n\n          const checkboxes = entries.map((e) => {\n            const label = e.folderPath ? `${e.folderPath}/${e.name}` : e.name;\n            return `<label><input type=\"checkbox\" class=\"s2s-entity-cb\" value=\"${e.id}\" /> ${label}</label>`;\n          }).join(\"\");\n\n          const selectAll = `<label class=\"s2s-select-all\"><input type=\"checkbox\" id=\"s2s-select-all\" /> Select All (${entries.length})</label>`;\n          listDiv.html(selectAll + checkboxes);\n\n          listDiv.find(\"#s2s-select-all\").on(\"change\", function () {\n            listDiv.find(\".s2s-entity-cb\").prop(\"checked\", this.checked);\n          });\n\n          updateTokensVisibility();\n        };\n\n        typeSelect.on(\"change\", loadEntities);\n        loadEntities();\n      },\n    }, { width: 480 }).render(true);\n  }\n\n  // --- Player Transfer Dialog ---\n  function showPlayerTransferDialog(targetId) {\n    const targetClient = clients.find((c) => c.id === targetId);\n    const targetName = targetClient?.customName || targetClient?.worldTitle || targetId;\n\n    // Get non-GM users from local game\n    const playerUsers = game.users.filter((u) => u.role < 4);\n    if (playerUsers.length === 0) {\n      ui.notifications.warn(\"No player accounts found.\");\n      return;\n    }\n\n    const userOptions = playerUsers\n      .map((u) => `<option value=\"${u.id}\">${u.name} (Role: ${u.role})</option>`)\n      .join(\"\");\n\n    const content = `\n      ${S2S_CSS}\n      <div class=\"s2s\">\n        <div class=\"s2s-header\">\n          <span>${sourceName}</span>\n          <i class=\"fas fa-arrow-right\"></i>\n          <span>${targetName}</span>\n        </div>\n        <div class=\"form-group\">\n          <label>Player</label>\n          <select name=\"userId\">${userOptions}</select>\n        </div>\n        <div class=\"s2s-row\">\n          <div class=\"form-group\">\n            <label><input type=\"checkbox\" name=\"createAccount\" checked /> Create account on target</label>\n          </div>\n          <div class=\"form-group\">\n            <label>Password</label>\n            <input type=\"text\" name=\"password\" placeholder=\"Optional\" />\n          </div>\n        </div>\n        <div class=\"s2s-checks\">\n          <label><input type=\"checkbox\" name=\"transferFiles\" checked /> Transfer files</label>\n          <label><input type=\"checkbox\" name=\"deleteSource\" /> Delete from source</label>\n          <label><input type=\"checkbox\" name=\"deleteTokens\" checked /> Delete tokens</label>\n          <label><input type=\"checkbox\" name=\"deleteUser\" /> Delete user account</label>\n        </div>\n      </div>\n    `;\n\n    new Dialog({\n      title: \"Transfer Player\",\n      content,\n      buttons: {\n        transfer: {\n          icon: '<i class=\"fas fa-exchange-alt\"></i>',\n          label: \"Transfer\",\n          callback: async (html) => {\n            const userId = html.find('[name=\"userId\"]').val();\n            const createAccount = html.find('[name=\"createAccount\"]').is(\":checked\");\n            const password = html.find('[name=\"password\"]').val();\n            const transferFilesOpt = html.find('[name=\"transferFiles\"]').is(\":checked\");\n            const deleteFromSource = html.find('[name=\"deleteSource\"]').is(\":checked\");\n            const deleteTokens = html.find('[name=\"deleteTokens\"]').is(\":checked\");\n            const deleteUser = html.find('[name=\"deleteUser\"]').is(\":checked\");\n\n            const user = game.users.get(userId);\n            if (!user) {\n              ui.notifications.error(\"User not found.\");\n              return;\n            }\n\n            showPlayerProgressDialog(user, targetId, {\n              createAccount,\n              password,\n              transferFilesOpt,\n              deleteFromSource,\n              deleteTokens,\n              deleteUser,\n            });\n          },\n        },\n        cancel: { icon: '<i class=\"fas fa-times\"></i>', label: \"Cancel\" },\n      },\n      default: \"transfer\",\n    }, { width: 450 }).render(true);\n  }\n\n  // --- Progress Dialog (entity transfer) ---\n  function showProgressDialog(entitiesToTransfer, targetId, options) {\n    const logLines = [];\n    const log = (msg) => {\n      logLines.push(msg);\n      const logEl = document.getElementById(\"s2s-log\");\n      if (logEl) {\n        logEl.textContent = logLines.join(\"\\n\");\n        logEl.scrollTop = logEl.scrollHeight;\n      }\n      const progEl = document.getElementById(\"s2s-progress\");\n      if (progEl) progEl.textContent = msg;\n    };\n\n    const content = `\n      ${S2S_CSS}\n      <div class=\"s2s\">\n        <p class=\"s2s-status\" id=\"s2s-progress\">Starting transfer...</p>\n        <pre class=\"s2s-log\" id=\"s2s-log\"></pre>\n      </div>\n    `;\n\n    new Dialog({\n      title: \"Transfer in Progress\",\n      content,\n      buttons: {\n        close: { icon: '<i class=\"fas fa-check\"></i>', label: \"Close\" },\n      },\n      default: \"close\",\n      render: async () => {\n        const results = [];\n\n        for (let i = 0; i < entitiesToTransfer.length; i++) {\n          const { entity, entityType } = entitiesToTransfer[i];\n          log(`--- Entity ${i + 1}/${entitiesToTransfer.length} ---`);\n          const result = await transferEntity(entity, entityType, targetId, options, log);\n          results.push(result);\n        }\n\n        const succeeded = results.filter((r) => r.success).length;\n        const failed = results.filter((r) => !r.success).length;\n\n        log(\"\\n========== TRANSFER COMPLETE ==========\");\n        log(`Entities: ${succeeded} succeeded, ${failed} failed`);\n\n        if (failed > 0) {\n          log(\"\\nFailed entities:\");\n          for (const r of results.filter((r) => !r.success)) {\n            log(`  - ${r.entityName}: ${r.error}`);\n          }\n        }\n\n        ui.notifications.info(`Transfer complete: ${succeeded} succeeded, ${failed} failed`);\n      },\n    }, { width: 520 }).render(true);\n  }\n\n  // --- Progress Dialog (player transfer) ---\n  function showPlayerProgressDialog(user, targetId, options) {\n    const logLines = [];\n    const log = (msg) => {\n      logLines.push(msg);\n      const logEl = document.getElementById(\"s2s-log\");\n      if (logEl) {\n        logEl.textContent = logLines.join(\"\\n\");\n        logEl.scrollTop = logEl.scrollHeight;\n      }\n    };\n\n    const content = `\n      ${S2S_CSS}\n      <div class=\"s2s\">\n        <p class=\"s2s-status\" id=\"s2s-progress\">Starting player transfer...</p>\n        <pre class=\"s2s-log\" id=\"s2s-log\"></pre>\n      </div>\n    `;\n\n    new Dialog({\n      title: \"Player Transfer in Progress\",\n      content,\n      buttons: {\n        close: { icon: '<i class=\"fas fa-check\"></i>', label: \"Close\" },\n      },\n      default: \"close\",\n      render: async () => {\n        try {\n          const results = await transferPlayer(user, targetId, options, log);\n          const succeeded = results.filter((r) => r.success).length;\n          const failed = results.filter((r) => !r.success).length;\n\n          log(\"\\n========== TRANSFER COMPLETE ==========\");\n          log(`Entities: ${succeeded} succeeded, ${failed} failed`);\n\n          if (failed > 0) {\n            log(\"\\nFailed entities:\");\n            for (const r of results.filter((r) => !r.success)) {\n              log(`  - ${r.entityName}: ${r.error}`);\n            }\n          }\n\n          ui.notifications.info(`Player transfer complete: ${succeeded} succeeded, ${failed} failed`);\n        } catch (err) {\n          log(`\\nFATAL ERROR: ${err.message}`);\n          ui.notifications.error(`Player transfer failed: ${err.message}`);\n        }\n      },\n    }, { width: 520 }).render(true);\n  }\n\n  // ========================= SECTION 8: MAIN ENTRY =========================\n\n  showServerSelectDialog();\n})();",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.j7NHQIWtXA4AG8Jt",
      "id": "j7NHQIWtXA4AG8Jt",
      "name": "Nuke",
      "type": "script",
      "author": "Gamemaster",
      "command": "async function cleanSlate() {\n    // List of document types to wipe\n    const collections = [\n        game.scenes,\n        game.actors,\n        game.items,\n        game.journal,\n        game.tables,\n        game.playlists,\n        game.cards,\n        // game.macros // Uncomment this line if you want to delete all macros too\n    ];\n\n    for (let collection of collections) {\n        const ids = collection.map(doc => doc.id);\n        if (ids.length > 0) {\n            console.log(`Deleting ${ids.length} documents from ${collection.name}...`);\n            await collection.documentClass.deleteDocuments(ids);\n        }\n    }\n\n    ui.notifications.info(\"World cleanup complete. A fresh start awaits!\");\n}\n\n// Confirmation Dialog\nnew Dialog({\n    title: \"Nuclear Option: Clear World Data\",\n    content: `\n        <div style=\"text-align: center;\">\n            <p><i class=\"fas fa-exclamation-triangle fa-3x\" style=\"color: #ff6b6b;\"></i></p>\n            <p>This will <strong>permanently delete</strong> all Scenes, Actors, Items, Journals, and more.</p>\n            <p><em>Are you absolutely sure?</em></p>\n        </div>`,\n    buttons: {\n        confirm: {\n            icon: '<i class=\"fas fa-trash\"></i>',\n            label: \"Delete Everything\",\n            callback: () => cleanSlate()\n        },\n        cancel: {\n            icon: '<i class=\"fas fa-times\"></i>',\n            label: \"Cancel\"\n        }\n    },\n    default: \"cancel\"\n}).render(true);",
      "img": "icons/svg/poison.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.AkcLmoRwvkrPvjyA",
      "id": "AkcLmoRwvkrPvjyA",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    }
  ]
}
```


