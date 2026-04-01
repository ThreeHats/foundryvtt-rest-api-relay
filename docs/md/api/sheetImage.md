---
tag: sheetimage
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# SheetImage

## GET /sheet/image

Get actor sheet as screenshot image

Captures the rendered actor sheet using html2canvas and returns it as a PNG or JPEG image. Works on both Foundry v12 and v13+.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string |  | query | The UUID of the entity to screenshot |
| selected | boolean |  | query | Whether to screenshot the selected entity's sheet |
| actor | boolean |  | query | Whether to use the selected token's actor if selected is true |
| clientId | string |  | query | Client ID for the Foundry world |
| format | string |  | query | Image format: png or jpeg (default: png) |
| quality | number |  | query | Image quality 0-1 for JPEG (default: 0.9) |
| scale | number |  | query | Capture scale factor (default: 1) |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**binary** - The sheet screenshot as an image

### Try It Out

<ApiTester
  method="GET"
  path="/sheet/image"
  parameters={[{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"quality","type":"number","required":false,"source":"query"},{"name":"scale","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

