# ---- Build Frontend ----
FROM node:20-alpine AS build
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install --legacy-peer-deps
COPY frontend ./frontend
RUN npm run build --prefix frontend

# ---- Runtime (Serve static) ----
FROM nginx:1.27-alpine
COPY --from=build /app/frontend/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
