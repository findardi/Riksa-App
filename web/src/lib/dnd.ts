// Drag payload types. A drag carrying `Files` is an upload from the OS; these
// two mark drags that started inside the app, and `dataTransfer.types` is the
// only part of a drag readable during `dragover` — so they are the switch.
export const FOLDER_MIME = 'application/x-riksa-folder';
export const DOCUMENT_MIME = 'application/x-riksa-document';
