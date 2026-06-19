// Workspace feature contract — mirrors internal/workspace/dto on the Go backend.

// Lifecycle state — server-validated enum; new rooms default to `prepare`.
export type WorkspaceStatus = 'prepare' | 'active' | 'archive';

// owner_id is NOT sent: the backend resolves it from the JWT claims.
export interface CreateWorkspacePayload {
	name: string;
	description: string;
}

// Backend PUT /workspaces/:id — name required, description optional.
export interface UpdateWorkspacePayload {
	name: string;
	description: string;
}

export interface WorkspaceData {
	id: string;
	owner_id: string;
	name: string;
	slug: string;
	description: string;
	status: WorkspaceStatus;
	// ISO-8601 strings over the wire (Go time.Time marshals to RFC3339).
	created_at: string;
	updated_at: string;
}

// Workspace Role
// Wire keys mirror the Go DTO: both create and update send `permissions` (plural).
export interface CreateWorkspaceRolePayload {
	permissions: string[];
	name: string;
	is_system: boolean;
}

export interface UpdateWorkspaceRolePayload {
	permissions: string[];
	name: string;
}

export interface WorkspaceRoleData {
	id: string;
	workspace_id: string;
	name: string;
	permissions: string[];
	is_system: boolean;
	created_at: string;
	updated_at: string;
}

// Workspace Member
export type MemberStatus = 'invited' | 'active' | 'suspended';

// Mirrors the Go `GetMemberResponse` (joined view: role name, user, groups).
export interface WorkspaceMemberData {
	id: string;
	workspace_id: string;
	user_id: string;
	role_id: string;
	status: MemberStatus;
	created_at: string;
	updated_at: string;
	role_name: string;
	username: string;
	email: string;
	// Go marshals a nil slice to null, so guard for null on the client.
	group_names: string[] | null;
}

export interface UpdateMemberRolePayload {
	role_id: string;
}
