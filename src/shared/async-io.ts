import { promisify } from 'util'
import * as fs from 'fs'
import { SpawnOptions, spawn } from 'child_process'
import * as gm from 'gm'

export const statAsync = promisify(fs.stat)
export const readFileAsync = promisify(fs.readFile)
export const writeFileAsync = promisify(fs.writeFile)
export const mkdirAsync = promisify(fs.mkdir)
export const readdirAsync = promisify(fs.readdir)
export const unlinkAsync = promisify(fs.unlink)
export const copyFileAsync = promisify(fs.copyFile)

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

export const fileExistsAsync = async (path: string, returnOnNotFile?: boolean): Promise<boolean> => {
  try {
    const stat = await statAsync(path)
    if (!stat.isFile()) {
      if (returnOnNotFile) {
        return false
      }
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

export type SpawnAsyncResult = {
  status: number;
  stdout: string;
  stderr: string;
}

export const spawnAsync = async (command: string, args: ReadonlyArray<string>, options?: SpawnOptions): Promise<SpawnAsyncResult> => {
  const child = spawn(command, args, options ? options : {})

  let stdout = ''
  let stderr = ''

  if (child.stdout) {
    for await (const chunk of child.stdout) {
      stdout += chunk
    }
  }

  if (child.stderr) {
    for await (const chunk of child.stderr) {
      stderr += chunk
    }
  }

  return new Promise(resolve => {
    child.on('close', code => resolve({
      status: code,
      stdout: stdout,
      stderr: stderr,
    }))
  })
}

export const writeGMAsync = async (filename: string, state: gm.State) => {
  await new Promise((resolve, reject) => {
    state.write(filename, err => {
      if (err) {
        reject(err)
      } else {
        resolve()
      }
    })
  })
}

export const getImageSize = async (filename: string): Promise<gm.Dimensions> => {
  return new Promise((resolve, reject) => {
    gm(filename).size((err, dims) => {
      if (err) {
        reject(err)
      } else {
        resolve(dims)
      }
    })
  })
}
