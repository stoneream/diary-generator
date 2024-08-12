package app

import com.typesafe.scalalogging.Logger

import java.io.File
import java.nio.file.{Files, Paths}

object DiaryArchiver {
  private val logger = Logger(DiaryArchiver.getClass)

  private def targetFinder(baseDirectory: File, startsWith: String): Array[File] = {
    require(baseDirectory.isDirectory)

    baseDirectory.listFiles.filter { file =>
      file.isDirectory && file.getName.startsWith(startsWith)
    }
  }

  private def completeDirectory(baseDirectory: File, startsWith: String): File = {
    require(baseDirectory.isDirectory)

    val archiveDir = baseDirectory.toPath.resolve("archive").resolve(startsWith).toFile

    if (archiveDir.exists() && archiveDir.isDirectory) {
      // do nothing
      logger.info(s"already exist : ${archiveDir.toString}")
    } else {
      logger.info(s"make directory : ${archiveDir.toString}")
      archiveDir.mkdirs()
    }

    archiveDir
  }

  def archive(baseDirectory: String, startsWith: String): Unit = {
    val baseDir = Paths.get(baseDirectory).toFile
    val archiveDir = completeDirectory(baseDir, startsWith)

    val moveTo = targetFinder(baseDir, startsWith).map { file =>
      (file.toPath, archiveDir.toPath.resolve(file.getName))
    }

    val infoText = moveTo
      .map { case (src, dest) =>
        s"${src} -> ${dest}"
      }
      .mkString("\n")

    logger.info(infoText)

    moveTo.foreach { case (src, dest) =>
      Files.move(src, dest)
    }
  }

}
