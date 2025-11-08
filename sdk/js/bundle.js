/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.
import { build } from 'esbuild'
import { glob } from 'glob'
import { spawn } from 'node:child_process'

function runCommand(command) {
  return new Promise((resolve, reject) => {
    const child = spawn(command, {
      stdio: 'inherit',
      shell: true,
    })
    child.on('close', (code) => {
      if (code !== 0) {
        reject(new Error(`Command failed with exit code ${code}`))
      } else {
        resolve()
      }
    })
    child.on('error', reject)
  })
}

async function buildProject() {
  try {
    // Find all TypeScript source files
    const entryPoints = await glob('src/**/*.ts', {
      ignore: ['src/**/*.spec.ts', 'src/**/*.test.ts'],
    })

    console.log(`Building ${entryPoints.length} modules...`)

    // Build without bundling to preserve module structure
    await build({
      entryPoints,
      bundle: false,
      outdir: 'dist/esm',
      format: 'esm',
      platform: 'neutral',
      target: 'es2020',
      sourcemap: true,
      logLevel: 'info',
      outExtension: { '.js': '.js' },
    })

    console.log('Generating type declarations...')
    await runCommand('tsc --emitDeclarationOnly --outDir dist/esm')

    console.log('✓ Build completed successfully')
  } catch (error) {
    console.error('Error during build:', error)
    process.exit(1)
  }
}

buildProject()
