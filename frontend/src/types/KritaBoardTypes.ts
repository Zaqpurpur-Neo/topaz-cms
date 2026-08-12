// Matches `CanvasItem` struct in Go
export interface CanvasItem {
  id: string;
  type: "main_artwork" | "reference";
  url: string;
  x: number; // Coords from translate(x, y) or 0
  y: number;
  width: number; // Render width from XML
  height: number; // Render height from XML
  opacity: number; // 0.0 to 1.0
  fileName: string;
}

// Matches `LayerNode` struct in Go (Recursive tree)
export interface LayerNode {
  name: string;
  nodeType: "paintlayer" | "grouplayer" | string;
  visible: boolean;
  opacity: number; // Normalized 0.0 to 1.0
  compositeOp: string; // "normal", "multiply", etc.
  filename: string; // Internal layer identifier (e.g., "layer3")
  children?: LayerNode[];
}

// Matches `Canvas` info in Go
export interface CanvasInfo {
  width: number;
  height: number;
  colorSpace: string;
  name: string;
}

// Matches `BoardResponse` struct returned by GET /api/board/:filename
export interface BoardResponse {
  fileName: string;
  canvas: CanvasInfo;
  items: CanvasItem[];
  layerTree: LayerNode[];
}
