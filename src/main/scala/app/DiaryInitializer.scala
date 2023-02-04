package app

import com.typesafe.scalalogging.Logger
import org.joda.time.DateTime

import java.io.File
import java.nio.charset.StandardCharsets
import java.nio.file.{Files, Paths}
import scala.io.Source
import scala.util.Using

object DiaryInitializer {
  private val logger = Logger(classOf[DiaryArchiver])

  /**
   * @param baseDirectory
   * @param templateFilePath
   * @param now
   * @return Either[ErrorReason, (TemplateFile, TargetFile)
   */
  def completeDirectory(baseDirectory: String, templateFilePath: String, now: DateTime): Either[String, (File, File)] = {
    val baseDir = Paths.get(baseDirectory)

    val templateFile = Paths.get(templateFilePath).toFile

    val targetDirPath = baseDir.resolve(Format.ymd(now)).toAbsolutePath
    val targetDir = targetDirPath.toFile
    val targetFilePath = targetDirPath.resolve(templateFile.getName)
    val targetFile = targetFilePath.toFile

    if (templateFile.exists()) {
      if (targetDir.exists()) {
        // do nothing
        logger.info(s"already exist : ${targetDirPath.toString}")
      } else {
        logger.info(s"make directory : ${targetDirPath.toString}")
        targetDir.mkdirs()
      }

      if (targetFile.exists()) {
        Left(s"already exist : ${targetFilePath.toString}")
      } else {
        logger.info(s"create file : ${targetFilePath.toString}")
        Right((templateFile, targetFile))
      }
    } else {
      Left("template file not found")
    }
  }

  def init(baseDirectory: String, templateFilePath: String, now: DateTime) = {
    completeDirectory(baseDirectory, templateFilePath, now) match {
      case Left(reason) => logger.error(reason)
      case Right((templateFile, targetFile)) =>
        Using(Source.fromFile(templateFile, "UTF-8", 4096)) { source => source.mkString }.fold(
          err => {
            logger.error("can't open template file", err)
          },
          template => {
            val formattedTemplate = TemplateLogic.render(template, now)
            targetFile.createNewFile()
            Files.write(targetFile.toPath, formattedTemplate.getBytes(StandardCharsets.UTF_8))
          }
        )
    }
  }
}
