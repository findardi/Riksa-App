// Presentation layer for the permission catalog. The *set* of valid permissions
// is owned by the Go backend (`permission.All`, fetched via getPermissions); here
// we only derive grouping from the `resource:action` naming and attach localized
// labels. Adding a permission on the backend surfaces automatically — no edit here.

import { t, type TKey } from '$lib/i18n';

// Canonical display order; anything unknown falls to the end, still grouped.
const RESOURCE_ORDER = ['workspace', 'member', 'role', 'group', 'folder', 'document'];
const ACTION_ORDER = ['view', 'create', 'add', 'upload', 'download', 'edit', 'assign', 'delete'];

const RESOURCE_LABEL: Record<string, TKey> = {
	workspace: 'perm.res.workspace',
	member: 'perm.res.member',
	role: 'perm.res.role',
	group: 'perm.res.group',
	folder: 'perm.res.folder',
	document: 'perm.res.document'
};

const ACTION_LABEL: Record<string, TKey> = {
	view: 'perm.act.view',
	create: 'perm.act.create',
	add: 'perm.act.add',
	upload: 'perm.act.upload',
	download: 'perm.act.download',
	edit: 'perm.act.edit',
	assign: 'perm.act.assign',
	delete: 'perm.act.delete'
};

// Seeded system roles carry technical names (owner/admin/…); show localized labels.
const SYSTEM_ROLE_LABEL: Record<string, TKey> = {
	owner: 'role.sys.owner',
	admin: 'role.sys.admin',
	contributor: 'role.sys.contributor',
	viewer: 'role.sys.viewer',
	guest: 'role.sys.guest'
};

// One-line capability summary for the fixed roles. Owner/admin/guest only —
// these are the roles the product commits to; anything else has no description.
const SYSTEM_ROLE_DESC: Record<string, TKey> = {
	owner: 'role.desc.owner',
	admin: 'role.desc.admin',
	guest: 'role.desc.guest'
};

/** Localized display name for a role; custom role names pass through unchanged. */
export function roleDisplayName(name: string): string {
	return SYSTEM_ROLE_LABEL[name] ? t(SYSTEM_ROLE_LABEL[name]) : name;
}

/** Localized one-line description for a fixed role; empty string if none. */
export function roleDescription(name: string): string {
	return SYSTEM_ROLE_DESC[name] ? t(SYSTEM_ROLE_DESC[name]) : '';
}

export type PermissionItem = { value: string; action: string; label: string };
export type PermissionGroup = { resource: string; label: string; items: PermissionItem[] };

function rank(list: string[], key: string): number {
	const i = list.indexOf(key);
	return i === -1 ? list.length : i;
}

function resourceLabel(resource: string): string {
	return RESOURCE_LABEL[resource] ? t(RESOURCE_LABEL[resource]) : resource;
}

function actionLabel(action: string): string {
	return ACTION_LABEL[action] ? t(ACTION_LABEL[action]) : action;
}

/** Group a flat catalog (e.g. `["document:view", ...]`) by resource, ordered. */
export function groupPermissions(all: string[]): PermissionGroup[] {
	const byResource = new Map<string, string[]>();
	for (const p of all) {
		const resource = p.split(':')[0] ?? p;
		const bucket = byResource.get(resource);
		if (bucket) bucket.push(p);
		else byResource.set(resource, [p]);
	}

	return [...byResource.entries()]
		.sort(([a], [b]) => rank(RESOURCE_ORDER, a) - rank(RESOURCE_ORDER, b))
		.map(([resource, values]) => ({
			resource,
			label: resourceLabel(resource),
			items: values
				.map((value) => {
					const action = value.split(':')[1] ?? value;
					return { value, action, label: actionLabel(action) };
				})
				.sort((a, b) => rank(ACTION_ORDER, a.action) - rank(ACTION_ORDER, b.action))
		}));
}
