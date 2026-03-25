---
tag: sheetImage
---

import ApiTester from '@site/src/components/ApiTester';

# sheetImage

## GET /sheet/image

Get actor sheet as screenshot image Captures the rendered actor sheet using html2canvas and returns it as a PNG or JPEG image. Works on both Foundry v12 and v13+.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| uuid | string |  | query | The UUID of the entity to screenshot |
| selected | boolean |  | query | Whether to screenshot the selected entity's sheet |
| actor | boolean |  | query | Whether to use the selected token's actor if selected is true |
| format | string |  | query | Image format: png or jpeg (default: png) |
| quality | number |  | query | Image quality 0-1 for JPEG (default: 0.9) |
| scale | number |  | query | Capture scale factor (default: 1) |
| userId | string |  | query | Foundry user ID or username to scope permissions |

### Returns

**binary** - The sheet screenshot as an image

### Try It Out

<ApiTester
  method="GET"
  path="/sheet/image"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"quality","type":"number","required":false,"source":"query"},{"name":"scale","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

