import type { UsersRecord } from "./pocketbase-types";

// type based on UserRecord but make the id and tokenKey optional
export type UserRequest = Omit<UsersRecord, "id" | "tokenKey"> & {
  passwordConfirm: string;
};

export type AuthStoreUserRecord = Omit<UsersRecord, "password" | "tokenKey">;

export interface MapInfo {
  tile_compression: string;
  tile_type: string;
  minzoom: number;
  maxzoom: number;
  bounds: number[];
  center: number[];
}

export interface MapSize {
  name: string;
  sizeBytes: number;
}

export interface GeoResponse {
  id: number;
  street: string;
  house_number: string;
  city: string;
  longitude: number;
  latitude: number;
}

export interface BackendInfo {
  reachable: boolean;
  name: string;
}
