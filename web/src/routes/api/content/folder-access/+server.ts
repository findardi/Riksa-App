import { error, json } from '@sveltejs/kit';
import { getFolderAccess, getFoldersTree } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { FolderAccessPanel, FolderTreeNode, InheritedFolderAccess } from '$lib/types/content';
import type { RequestHandler } from './$types';

function trail(
	nodes: FolderTreeNode[],
	id: string,
	path: FolderTreeNode[] = []
): FolderTreeNode[] | null {
	for (const n of nodes) {
		const next = [...path, n];
		if (n.id === id) return next;
		const hit = trail(n.children ?? [], id, next);
		if (hit) return hit;
	}
	return null;
}

export const GET: RequestHandler = async ({ locals, url }) => {
	const token = locals.session;
	if (!token) error(401, t('err.invalidCredentials'));

	const workspaceId = url.searchParams.get('workspaceId');
	const folderId = url.searchParams.get('folderId');
	if (!workspaceId || !folderId) error(400, t('err.generic'));

	const tree = await getFoldersTree(token, workspaceId);
	if (!tree.ok) error(tree.status || 500, tree.message);

	const path = trail(tree.data, folderId);
	if (!path) error(404, t('facc.err.notFound'));

	const ancestors = path.slice(0, -1).reverse();

	const [self, ...chain] = await Promise.all([
		getFolderAccess(token, workspaceId, folderId),
		...ancestors.map((a) => getFolderAccess(token, workspaceId, a.id))
	]);

	if (!self.ok) error(self.status || 500, self.message);

	const direct = self.data;
	const claimed = new Set(direct.map((r) => r.group_id));
	const inherited: InheritedFolderAccess[] = [];

	ancestors.forEach((ancestor, i) => {
		const res = chain[i];
		if (!res?.ok) return;
		for (const row of res.data) {
			if (claimed.has(row.group_id)) continue;
			claimed.add(row.group_id);
			inherited.push({
				...row,
				source_folder_id: ancestor.id,
				source_folder_name: ancestor.name
			});
		}
	});

	const panel: FolderAccessPanel = { folder_id: folderId, direct, inherited };
	return json(panel);
};
