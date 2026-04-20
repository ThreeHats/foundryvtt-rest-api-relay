/** Human-readable descriptions for each scope, shown in approval pages. */
export const scopeDescriptions: Record<string, string> = {
  'entity:read':      'Read actors, items, journal entries, and other game entities',
  'entity:write':     'Create and modify actors, items, journal entries, and other game entities',
  'roll:read':        'Read recent roll history',
  'roll:execute':     'Execute dice rolls',
  'chat:read':        'Read chat messages',
  'chat:write':       'Send chat messages',
  'encounter:read':   'Read combat encounters and initiative order',
  'encounter:manage': 'Create and manage combat encounters',
  'scene:read':       'Read scene data and canvas state',
  'scene:write':      'Create and modify scenes',
  'canvas:read':      'Read canvas token positions and placeables',
  'canvas:write':     'Move and modify tokens and other canvas placeables',
  'effects:read':     'Read active effects on tokens and actors',
  'effects:write':    'Apply and remove active effects',
  'macro:list':       'List available macros',
  'macro:execute':    'Execute Foundry macros',
  'macro:write':      'Create and modify Foundry macros',
  'user:read':        'Read user and player information',
  'user:write':       'Modify user settings',
  'file:read':        'Read files from the Foundry server',
  'file:write':       'Upload files to the Foundry server',
  'structure:read':   'Read compendium and folder structure',
  'structure:write':  'Modify compendium and folder structure',
  'sheet:read':       'Open and interact with character sheets',
  'playlist:control': 'Control music playlists',
  'world:info':       'Read world metadata and settings',
  'clients:read':     'List connected clients and sessions',
  'events:subscribe': 'Subscribe to real-time game events via SSE',
  'search':           'Search game entities using QuickInsert',
  'session:manage':   'Start and stop headless browser sessions for automated GM actions',
  'execute-js':       'Execute arbitrary JavaScript in the Foundry game session',
  'dnd5e':            'Access D&D 5e system-specific endpoints',
};

/**
 * Danger details for scopes that warrant an expanded warning.
 * Shown as an alert box in addition to the normal description.
 */
export const dangerDetails: Record<string, string> = {
  'execute-js':    'Can read or modify any game data, access server settings, and potentially exfiltrate data. Only grant to applications you fully trust.',
  'macro:execute': 'Macros may contain arbitrary JavaScript. Only grant if you trust the requesting application not to run harmful macros.',
  'macro:write':   'Can create macros with arbitrary code that could be executed later. Only grant to fully trusted applications.',
};
