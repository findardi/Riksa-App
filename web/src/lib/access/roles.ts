// Role-based capability gate. The backend (`permission` pkg) is the real
// authority; this mirrors the *intent* of the three system roles so the UI can
// hide controls and the room layout can block routes. Server load still enforces
// the gate — these helpers are convenience, not security.
//
// Strict three-role model: owner / admin / guest. Any other role name (e.g. a
// legacy contributor/viewer still returned by the backend) is treated as guest,
// i.e. least privilege.

export type WorkspaceRole = 'owner' | 'admin' | 'guest';

export function normalizeRole(role: string): WorkspaceRole {
	return role === 'owner' || role === 'admin' ? role : 'guest';
}

const MANAGER: WorkspaceRole[] = ['owner', 'admin'];

/** May open the access-management area (members / roles / groups). */
export function canManageAccess(role: string): boolean {
	return MANAGER.includes(normalizeRole(role));
}

/** May invite, change role, or remove members. */
export function canManageMembers(role: string): boolean {
	return MANAGER.includes(normalizeRole(role));
}

/** May create, edit, or delete workspace roles. */
export function canManageRoles(role: string): boolean {
	return MANAGER.includes(normalizeRole(role));
}

/** May create, edit, or delete groups. */
export function canManageGroups(role: string): boolean {
	return MANAGER.includes(normalizeRole(role));
}

// Workspace-level mutations (edit details, change status, delete) are all
// guarded by `RequireOwner` on the backend (workspace/module.go) — owner only,
// NOT admin. Keep these two split in case the backend later relaxes edit/status.

/** Owner-only: edit workspace details and change its status. */
export function canEditWorkspace(role: string): boolean {
	return normalizeRole(role) === 'owner';
}

/** Owner-only: delete the workspace itself. */
export function canDeleteWorkspace(role: string): boolean {
	return normalizeRole(role) === 'owner';
}
