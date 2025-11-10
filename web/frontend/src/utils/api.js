import { GEO_HOST } from '../config';

const withOrgId = (url) => {
  const orgId = localStorage.getItem("selectedOrgId");
  if (!orgId) {
    return url;
  }
  const separator = url.includes('?') ? '&' : '?';
  return `${url}${separator}org_id=${orgId}`;
};

const getHeaders = (hasBody = false) => {
  const headers = {};
  if (hasBody) {
    headers["Content-Type"] = "application/json";
  }
  return headers;
};

export const apiClient = {
  get: async (url) => {
    const response = await fetch(withOrgId(url), {
      method: "GET",
      headers: getHeaders(),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  post: async (url, data) => {
    const response = await fetch(withOrgId(url), {
      method: "POST",
      headers: getHeaders(true),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  put: async (url, data) => {
    const response = await fetch(withOrgId(url), {
      method: "PUT",
      headers: getHeaders(true),
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
  delete: async (url) => {
    const response = await fetch(withOrgId(url), {
      method: "DELETE",
      headers: getHeaders(),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }));
      throw new Error(errorData.message || 'Request failed');
    }
    return response.json();
  },
};

export const getDeviceHistory = async (mapId, deviceId, from, to, take, skip) => {
  const url = `${GEO_HOST}/api/v1/geo/history?map_id=${mapId}&from=${from}&to=${to}&take=${take}&total=true&skip=${skip}&entity_id=${deviceId}`;
  return apiClient.get(url);
};
