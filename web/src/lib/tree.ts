import type { FolderTreeNode } from '$lib/types/content';

export function findNode(nodes: FolderTreeNode[], id: string): FolderTreeNode | null {
	for (const n of nodes) {
		if (n.id === id) return n;
		const hit = findNode(n.children ?? [], id);
		if (hit) return hit;
	}
	return null;
}
