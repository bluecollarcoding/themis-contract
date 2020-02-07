import { promisify } from 'util'
import * as fs from 'fs'

export const statAsync = promisify(fs.stat)
export const readFileAsync = promisify(fs.readFile)
export const writeFileAsync = promisify(fs.writeFile)
export const mkdirAsync = promisify(fs.mkdir)
export const readdirAsync = promisify(fs.readdir)
export const unlinkAsync = promisify(fs.unlink)

export class IOError extends Error { }

export const dirExistsAsync = async (path: string): Promise<boolean> => {
  try {
    const stat = await statAsync(path)
    if (!stat.isDirectory()) {
      throw new IOError(`Path exists but is not a directory: ${path}`)
    }
    return true
  } catch (error) {
    if (error instanceof IOError) {
      throw error
    }
    return false
  }
}

/**
 * Ensures that the given path is a directory, creating parent directories if
 * necessary.
 * @param {string} path The path to ensure exists.
 */
export const ensurePath = async (path: string) => {
  if (!(await dirExistsAsync(path))) {
    await mkdirAsync(path, { recursive: true })
  }
}

export const fileExistsAsync = async (path: string): Promise<boolean> => {
  try {
    const stat = await statAsync(path)
    if (!stat.isFile()) {
      throw new IOError(`Path exists but is not a file: ${path}`)
    }
    return true
  } catch (error) {
    if (error instanceof IOError) {
      throw error
    }
    return false
  }
}
